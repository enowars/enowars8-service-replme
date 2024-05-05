import { Terminal } from "@xterm/xterm";
import { AttachAddon } from "@xterm/addon-attach";
import { CanvasAddon } from "@xterm/addon-canvas";
import { ClipboardAddon } from "@xterm/addon-clipboard";
import { FitAddon } from "@xterm/addon-fit";

const retry = 5;
const delay = 1000;

let pid: string | undefined;
const protocol = (location.protocol === 'https:') ? 'wss://' : 'ws://';

const urlParams = new URLSearchParams(window.location.search);
const port = urlParams.get('port');

let serverURL = location.protocol + '//' + location.hostname + ':' + port
let socketURL = protocol + location.hostname + ':' + port + '/terminals/';


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
  const url = serverURL + '/terminals/' + pid + '/size?cols=' + cols + '&rows=' + rows;

  fetch(url, { method: 'POST' });
});

const resizeObserver = new ResizeObserver(() => {
  fit.fit();
});

fit.fit()

let terminalContainer = document.getElementById('terminal');
terminal.open(terminalContainer!);
terminal.focus();

resizeObserver.observe(terminalContainer!);

async function sleep(ms: number) {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve(null);
    }, ms);
  })
}

async function connect() {
  for (let i = 0; i < retry; i++) {
    try {
      const res = await fetch(
        serverURL + '/terminals?cols=' + terminal.cols + '&rows=' + terminal.rows,
        { method: 'POST' }
      )
      if (res.ok) {
        pid = await res.text();
        socketURL += pid;
        const socket = new WebSocket(socketURL);
        socket.onopen = () => {
          terminal.loadAddon(new AttachAddon(socket));
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
