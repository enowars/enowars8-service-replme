import { RequestHandler } from "express";
import { WebsocketRequestHandler } from "express-ws";

const express: RequestHandler = (req, res, next) => {
  const key = req.params.key;
  if (key === undefined || typeof (key) !== 'string') {
    res.status(401);
    res.send({ error: "Unauthorized" });
    return;
  }

  if (key !== process.env.API_KEY) {
    res.status(401);
    res.send({ error: "API-Key invalid" });
    return;
  }

  next();
}

const ws: WebsocketRequestHandler = (ws, req, next) => {
  const key = req.params.key;
  if (key === undefined || typeof (key) !== 'string') {
    ws.close();
    return;
  }

  if (key !== process.env.API_KEY) {
    ws.close();
    return;
  }

  next();
}

const ApiMiddleware = {
  express,
  ws
}

export default ApiMiddleware;
