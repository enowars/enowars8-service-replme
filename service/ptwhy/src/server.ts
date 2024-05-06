import fs from "node:fs";
import child_process from "node:child_process";

import express from "express";
import expressWs from "express-ws";
import WebSocket from "ws";
import * as pty from "node-pty";
import cors from "cors";

import z from "zod";

const CreateUserScheme = z.object({
  username: z.string().min(4).max(64).regex(/^[a-zA-Z0-9]*$/),
  password: z.string().min(4).max(64).regex(/^[a-zA-Z0-9]*$/)
});

function main() {
  const app = express();
  app.use(cors())
  app.use(express.json())

  const appWs = expressWs(app).app;

  const terminals = {};
  const unsentOutput = {};
  const temporaryDisposable = {};

  app.post('/user/create', (req, res) => {

    let schema = CreateUserScheme.safeParse(req.body);
    if (!schema.success) {
      res.status(400)
      res.send({ error: "Invalid input" });
      return;
    }

    console.log("here")

    let body = schema.data;

    const shadow = fs.readFileSync('/etc/shadow', 'utf8').split('\n');

    const user = shadow.find((user) => {
      return user.startsWith(body.username)
    })

    if (user) {
      const [, hash] = user.split(':')
      const _hash = hash.split('$')
      if (_hash.length < 4) {
        res.status(403)
        res.send({ error: "Forbidden" })
        return;
      }
      let type: string, salt: string, pw: string = '';
      if (hash[1] === 'y') {
        res.status(403)
        res.send({ error: "Forbidden" })
        return;
      } else {
        [, type, salt, pw] = _hash;
      }

      const buf = child_process.execSync(`openssl passwd -${type} -salt ${salt} ${body.password}`)
      const result = buf.toString('utf8').trim()
      if (result === hash) {
        res.status(200);
        res.send({ success: "Ok" });
        return;
      } else {
        res.status(401);
        res.send({ error: "Invalid credentials" });
        return;
      }
    }

    child_process.execSync(`adduser -D ${body.username} -s /bin/zsh`);
    child_process.execSync(`echo "${body.username}:${body.password}" | chpasswd`);

    res.status(201);
    res.send({ success: "Created" });
  });

  app.post('/terminals', (req, res) => {
    console.log("terminals");
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
  });

  app.post('/terminals/:pid/size', (req, res) => {
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
  });

  appWs.ws('/terminals/:pid', function (ws, req) {
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
    ws.on('close', function () {
      term.kill();
      console.log('Closed terminal ' + term.pid);
      delete terminals[term.pid];
    });
  });

  const port = 3000;
  const host = '0.0.0.0';

  console.log('App listening to http://' + host + ':' + port);
  app.listen(port, host, 0);
}

main();
