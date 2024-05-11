import express from "express";
import expressWs from "express-ws";
import cors from "cors";

import UserController from "controller/user-controller";
import TermController from "controller/term-controller";

export default function createExpressApp() {
  const app = express();
  app.use(cors())
  app.use(express.json())

  const appWs = expressWs(app).app;

  app.post('/user/login', UserController.login);
  app.post('/user/create', UserController.create);
  app.post('/terminals', TermController.create);
  app.post('/terminals/:pid/size', TermController.resize);
  appWs.ws('/terminals/:pid', TermController.websocket);

  return app;
}

