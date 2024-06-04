import { RequestHandler } from "express";
import { WebsocketRequestHandler } from "express-ws";
import WebSocket from "ws";
import * as pty from "node-pty";

const terminals = {};
const unsentOutput = {};
const temporaryDisposable = {};

const create: RequestHandler = (req, res) => {
  const env: { [key: string]: string } = {};
  for (const k of Object.keys(process.env)) {
    const v = process.env[k];
    if (v) {
      env[k] = v;
    }
  }
  env['COLORTERM'] = 'truecolor';
  if (typeof req.query.cols !== 'string' || typeof req.query.rows !== 'string') {
    console.error({ req });
    throw new Error('Unexpected query args');
  }
  const cols = parseInt(req.query.cols);
  const rows = parseInt(req.query.rows);
  const term = pty.spawn('login', [], {
    name: 'xterm-256color',
    cols: cols ?? 80,
    rows: rows ?? 24,
    cwd: env.PWD,
    env,
    encoding: null,
  });

  console.log('Created terminal with PID: ' + term.pid);
  terminals[term.pid] = term;
  unsentOutput[term.pid] = '';
  temporaryDisposable[term.pid] = term.onData(function(data) {
    unsentOutput[term.pid] += data;
  });
  res.send(term.pid.toString());
  res.end();
}

const resize: RequestHandler = (req, res) => {
  if (typeof req.query.cols !== 'string' || typeof req.query.rows !== 'string') {
    console.error({ req });
    throw new Error('Unexpected query args');
  }
  const pid = parseInt(req.params.pid);
  const cols = parseInt(req.query.cols);
  const rows = parseInt(req.query.rows);
  const term = terminals[pid];

  term.resize(cols, rows);
  console.log('Resized terminal ' + pid + ' to ' + cols + ' cols and ' + rows + ' rows.');
  res.end();
}

const websocket: WebsocketRequestHandler = (ws, req) => {
  const term = terminals[parseInt(req.params.pid)];
  console.log('Connected to terminal ' + term.pid);
  temporaryDisposable[term.pid].dispose();
  delete temporaryDisposable[term.pid];
  ws.send(unsentOutput[term.pid]);
  delete unsentOutput[term.pid];

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

  term.onData(function(data: any) {
    try {
      send(data);
    } catch (error) {
      // websocket is not open, ignore
    }
  });
  ws.on('message', function(msg) {
    term.write(msg);
    userInput = true;
  });
  ws.on('close', function() {
    term.kill('SIGKILL');
    console.log('Closed terminal ' + term.pid);
    delete terminals[term.pid];
  });
}

const TermController = {
  create,
  resize,
  websocket
}

export default TermController;

