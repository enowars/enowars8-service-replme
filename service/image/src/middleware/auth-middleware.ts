import { RequestHandler } from "express";
import { WebsocketRequestHandler } from "express-ws";

const express: RequestHandler = (req, res, next) => {
  const user = req.session.user
  if (!user) {
    res.status(401);
    res.send({ error: "Unauthorized" });
    return;
  }
  next();
}

const ws: WebsocketRequestHandler = (ws, req, next) => {
  const user = req.session.user
  if (!user) {
    ws.close()
    return;
  }
  next();
}

const AuthMiddleware = {
  express,
  ws
}

export default AuthMiddleware;
