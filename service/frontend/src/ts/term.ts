import { Terminal } from "@xterm/xterm";
import { AttachAddon } from "@xterm/addon-attach";
import { CanvasAddon } from "@xterm/addon-canvas";
import { ClipboardAddon } from "@xterm/addon-clipboard";
import { FitAddon } from "@xterm/addon-fit";
import { getPassword, sleep } from "./lib";

const retry = 5;
const delay = 1000;

let pid: string | undefined;
const protocol = (location.protocol === 'https:') ? 'wss://' : 'ws://';

const [username, port] = location.pathname.replace('/term/', '').split('/');
const password = getPassword(username)

let socketURL = protocol + location.host + '/ws/terminal/' + port + '/';

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

terminal.onResize((size: { cols: number, rows: number }) => {
  if (!pid) {
    return;
  }
  const cols = size.cols;
  const rows = size.rows;

  fetch(
    `/api/ptwhy/terminals/${pid}/size`,
    {
      method: 'POST',
      headers: {
        'Authorization': 'Basic ' + btoa(username + ":" + password),
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        cols,
        rows
      })
    }
  );
});

const resizeObserver = new ResizeObserver(() => {
  fit.fit();
});

fit.fit()

let terminalContainer = document.getElementById('terminal');
terminal.open(terminalContainer!);
terminal.focus();

resizeObserver.observe(terminalContainer!);

async function connect() {
  for (let i = 0; i < retry; i++) {
    try {
      const res = await fetch(
        '/api/ptwhy/terminals',
        {
          method: 'POST',
          headers: {
            'Authorization': 'Basic ' + btoa(username + ":" + password),
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            cols: terminal.cols,
            rows: terminal.rows
          })
        }
      )
      if (res.ok) {
        pid = await res.text();
        socketURL += pid;
        const socket = new WebSocket(socketURL);
        socket.onopen = async () => {
          await sleep(100);
          socket.send(`${username}\n`);
          await sleep(100);
          socket.send(`${password}\n`);
          terminal.loadAddon(new AttachAddon(socket));
        };
        window.onbeforeunload = function (e: any) {
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
