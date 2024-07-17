"use client";

import { useEffect, useRef } from "react";
import { Terminal as XTerm } from "@xterm/xterm";
import { CanvasAddon } from "@xterm/addon-canvas";
import { ClipboardAddon } from "@xterm/addon-clipboard";
import { FitAddon } from "@xterm/addon-fit";
import { AttachAddon } from "@/lib/attach-addon";
import { sleep } from "@/lib/utils";
import { useTheme } from "next-themes";

const LIGHT_THEME = {
  background: "#ffffff",
  foreground: "#333333",
  cursor: "#333333",
  cursorAccent: "#ffffff",
  selectionBackground: "#add6ff",
  black: "#000000",
  blue: "#0451a5",
  brightBlack: "#666666",
  brightBlue: "#0451a5",
  brightCyan: "#0598bc",
  brightGreen: "#14ce14",
  brightMagenta: "#bc05bc",
  brightRed: "#cd3131",
  brightWhite: "#a5a5a5",
  brightYellow: "#b5ba00",
  cyan: "#0598bc",
  green: "#00bc00",
  magenta: "#bc05bc",
  red: "#cd3131",
  white: "#555555",
  yellow: "#949800",
};

const DARK_THEME = {
  foreground: "#f8f8f8",
  background: "#020817",
  selectionBackground: "#5da5d533",
  selectionInactiveBackground: "#555555aa",
  black: "#1e1e1d",
  brightBlack: "#262625",
  red: "#ce5c5c",
  brightRed: "#ff7272",
  green: "#5bcc5b",
  brightGreen: "#72ff72",
  yellow: "#cccc5b",
  brightYellow: "#ffff72",
  blue: "#5d5dd3",
  brightBlue: "#7279ff",
  magenta: "#bc5ed1",
  brightMagenta: "#e572ff",
  cyan: "#5da5d5",
  brightCyan: "#72f0ff",
  white: "#f8f8f8",
  brightWhite: "#ffffff",
};

const RETRY = 5;
const DELAY = 1000;

type TerminalProps = {
  id?: string;
  className?: string;
  path: string;
  catchClose?: boolean;
};

const Terminal: React.FC<TerminalProps> = (props) => {
  const { id, className, path, catchClose } = props;

  const { resolvedTheme } = useTheme();

  const terminalRef = useRef<HTMLDivElement>(null);
  const xtermRef = useRef<XTerm | null>(null);
  const websocketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!xtermRef.current) return;

    xtermRef.current.options.theme = {
      ...(resolvedTheme === "light" ? LIGHT_THEME : DARK_THEME),
    };
  }, [xtermRef, resolvedTheme]);

  useEffect(() => {
    if (!terminalRef.current) return;

    const protocol = location.protocol === "https:" ? "wss://" : "ws://";
    let socketURL =
      protocol + (process.env.NEXT_PUBLIC_WS || location.host) + path;
    const terminal = new XTerm({
      allowProposedApi: true,
      fontFamily: '"DejaVuSansM Nerd Font", courier-new, courier, monospace',
      theme: resolvedTheme === "light" ? LIGHT_THEME : DARK_THEME,
    });

    xtermRef.current = terminal;

    const canvas = new CanvasAddon();
    const clipboard = new ClipboardAddon();
    const fit = new FitAddon();

    terminal.loadAddon(canvas);
    terminal.loadAddon(clipboard);
    terminal.loadAddon(fit);

    const resizeObserver = new ResizeObserver(() => {
      fit.fit();
    });

    terminal.open(terminalRef.current);
    terminal.focus();
    resizeObserver.observe(terminalRef.current);

    async function connect() {
      for (let i = 0; i < RETRY; i++) {
        try {
          const socket = new WebSocket(socketURL);
          websocketRef.current = socket;
          socket.onopen = async () => {
            terminal.loadAddon(new AttachAddon(socket));
            fit.fit();
          };
          if (catchClose)
            window.onbeforeunload = function (e: any) {
              if (e) {
                e.returnValue = "Leave site?";
              }
              // safari
              return "Leave site?";
            };
          break;
        } catch (error) {
          // ignore
        } finally {
          await sleep(DELAY);
        }
      }
    }

    connect();

    return () => {
      terminal.dispose();
      websocketRef.current?.close();
      if (catchClose) window.onbeforeunload = null;
    };
  }, [terminalRef, websocketRef]);

  return <div id={id} className={className} ref={terminalRef} />;
};

export default Terminal;
