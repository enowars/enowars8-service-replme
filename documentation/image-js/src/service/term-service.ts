import * as pty from "node-pty";
import WebSocket from "ws";
import { randomUUID } from "crypto";
import UserService from "./user-service";
import { TermServiceError } from "../types/error";

const terminals: Record<string, pty.IPty> = {};
const unsentOutput: Record<string, string> = {};
const temporaryDisposable: Record<string, pty.IDisposable> = {};

type CreateOptions = {
  executable?: string,
  cols?: number,
  rows?: number
};

type ResizeOptions = {
  cols: number,
  rows: number
}

function create(username: string, opts?: CreateOptions) {
  const env: Record<string, string> = {};
  env['COLORTERM'] = 'truecolor';
  env['PATH'] = process.env.PATH;

  const userinfo = UserService.info(username);

  const term = pty.spawn(opts?.executable ?? userinfo.shell, [], {
    name: 'xterm-256color',
    cols: opts?.cols ?? 80,
    rows: opts?.rows ?? 24,
    uid: userinfo.uid,
    gid: userinfo.gid,
    cwd: userinfo.home,
    env,
    encoding: null,
  });

  const uuid = randomUUID();

  terminals[uuid] = term;
  unsentOutput[uuid] = '';
  temporaryDisposable[uuid] = term.onData(function(data) {
    unsentOutput[uuid] += data;
  });

  return uuid;
}

function resize(uuid: string, opts: ResizeOptions) {
  const term = terminals[uuid];
  if (!term)
    throw new TermServiceError(404, "term not found");
  term.resize(opts.cols, opts.rows);
}

function ws(uuid: string, ws: WebSocket, onClose?: () => void) {
  const term = terminals[uuid];
  if (!term)
    throw new TermServiceError(404, "term not found");

  temporaryDisposable[uuid].dispose();
  delete temporaryDisposable[uuid];
  ws.send(unsentOutput[uuid]);
  delete unsentOutput[uuid];

  let userInput = false;

  function bufferUtf8(socket: WebSocket, timeout: number, maxSize: number) {
    const chunks = [];
    let length = 0;
    let sender = null;
    return (data: any) => {
      chunks.push(data);
      length += data.length;
      if (length > maxSize || userInput) {
        userInput = false;
        socket.send(Buffer.concat(chunks));
        chunks.length = 0;
        length = 0;
        if (sender) {
          clearTimeout(sender);
          sender = null;
        }
      } else if (!sender) {
        sender = setTimeout(() => {
          socket.send(Buffer.concat(chunks));
          chunks.length = 0;
          length = 0;
          sender = null;
        }, timeout);
      }
    };
  }
  const send = bufferUtf8(ws, 3, 262144);

  term.onData(function(data: string) {
    try {
      send(data);
    } catch (error) {
      // websocket is not open, ignore
    }
  });
  ws.on('message', function(msg: string) {
    term.write(msg);
    userInput = true;
  });
  ws.on('close', function() {
    term.kill('SIGKILL');
    delete terminals[uuid];
    if (onClose)
      onClose();
  });

}

const TermService = {
  create,
  resize,
  ws
}

export default TermService;
