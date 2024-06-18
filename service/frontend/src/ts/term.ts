import { Terminal } from "@xterm/xterm";
import { CanvasAddon } from "@xterm/addon-canvas";
import { ClipboardAddon } from "@xterm/addon-clipboard";
import { FitAddon } from "@xterm/addon-fit";
import { AttachAddon, randomString, sleep } from "./lib";

const retry = 5;
const delay = 1000;

let replUuid: string | undefined;
const protocol = (location.protocol === 'https:') ? 'wss://' : 'ws://';

const [__1, __2, name] = location.pathname.split('/');

let socketURL = protocol + location.host + '/api/repl/';

const terminal = new Terminal({
  allowProposedApi: true,
  fontFamily: '"DejaVuSansM Nerd Font", courier-new, courier, monospace',
  theme: {
    background: '#ffffff',
    foreground: '#333333',
    cursor: '#333333',
    cursorAccent: '#ffffff',
    selectionBackground: '#add6ff',
    black: '#000000',
    blue: '#0451a5',
    brightBlack: '#666666',
    brightBlue: '#0451a5',
    brightCyan: '#0598bc',
    brightGreen: '#14ce14',
    brightMagenta: '#bc05bc',
    brightRed: '#cd3131',
    brightWhite: '#a5a5a5',
    brightYellow: '#b5ba00',
    cyan: '#0598bc',
    green: '#00bc00',
    magenta: '#bc05bc',
    red: '#cd3131',
    white: '#555555',
    yellow: '#949800'
  }
});

const canvas = new CanvasAddon();
const clipboard = new ClipboardAddon();
const fit = new FitAddon();

terminal.loadAddon(canvas);
terminal.loadAddon(clipboard);
terminal.loadAddon(fit);

const resizeObserver = new ResizeObserver(() => {
  fit.fit();
});

let terminalContainer = document.getElementById('terminal');
terminal.open(terminalContainer!);
terminal.focus();

resizeObserver.observe(terminalContainer!);

async function connect() {
  for (let i = 0; i < retry; i++) {
    try {
      let res: Response;
      if (name) {
        res = await fetch(
          '/api/repl/name/' + name,
          {
            method: 'POST',
          }
        )
      } else {
        const username = randomString(60);
        const password = randomString(60);
        res = await fetch(
          '/api/repl',
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              username,
              password
            })
          }
        )
      }
      if (res.ok) {
        replUuid = (await res.json()).replUuid;
        socketURL += replUuid;
        const socket = new WebSocket(socketURL);
        socket.onopen = async () => {
          document.addEventListener('keydown', (event) => {
            if (event.code === "Escape") {
              event.preventDefault();
              socket.send(JSON.stringify({ stdin: '\x1B' }));
              terminal.focus();
            }
          });
          terminal.loadAddon(new AttachAddon(socket));
          fit.fit()
        };
        window.onbeforeunload = function(e: any) {
          if (e) {
            e.returnValue = 'Leave site?';
          }
          // safari
          return 'Leave site?';
        };
        break;
      }
    } catch (error) {
      // ignore
    } finally {
      await sleep(delay);
    }
  }
}

connect();
