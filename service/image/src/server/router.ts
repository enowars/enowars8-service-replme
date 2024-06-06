import express from "express";
import expressWs from "express-ws";
import cors from "cors";
import session from "express-session";
import memorystore from "memorystore";
import crypto from "crypto";
import morgan from "morgan";

import AuthController from "../controller/auth-controller";
import TermController from "../controller/term-controller";
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

  app.use('/api/:key', ApiMiddleware.express);
  app.post('/api/:key/auth/register', AuthController.register);
  app.post('/api/:key/auth/login', AuthController.login);

  app.use('/api/:key/term', AuthMiddleware.express)

  app.post('/api/:key/term', TermController.create);
  appWs.ws('/api/:key/term', TermController.ws);
  app.post('/api/:key/term/size', TermController.resize);

  return app;
}

