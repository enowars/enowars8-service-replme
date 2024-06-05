import express from "express";
import expressWs from "express-ws";
import cors from "cors";
import session from "express-session";
import memorystore from "memorystore";
import crypto from "crypto";
import morgan from "morgan";

import AuthController from "../controller/auth-controller";
import UserController from "../controller/user-controller";
import TermController from "../controller/term-controller";
import TermV2Controller from "../controller/term-v2-controller";
import AuthMiddleware from "../middleware/auth-middleware";
import ApiMiddleware from "../middleware/api-middleware";

export default function createExpressApp() {
  const app = express();

  const MemoryStore = memorystore(session)

  app.use(session({
    cookie: { maxAge: 86400000 },
    store: new MemoryStore({
      checkPeriod: 86400000
    }),
    resave: false,
    secret: crypto.randomBytes(64).toString()
  }));

  app.use(cors());
  app.use(express.json());
  app.use(morgan('combined'));

  const appWs = expressWs(app).app;

  app.post('/user/login', UserController.login);
  app.post('/user/create', UserController.create);
  app.post('/terminals', TermController.create);
  app.post('/terminals/:pid/size', TermController.resize);
  appWs.ws('/terminals/:pid', TermController.websocket);

  ////////////////////// APIV2 //////////////////////

  app.use('/api/:key', ApiMiddleware.express);
  app.post('/api/:key/auth/register', AuthController.register);
  app.post('/api/:key/auth/login', AuthController.login);

  app.use('/api/:key/term', AuthMiddleware.express)
  app.post('/api/:key/term', TermV2Controller.create);
  appWs.ws('/api/:key/term', TermV2Controller.ws);
  app.post('/api/:key/term/size', TermV2Controller.resize);

  return app;
}

