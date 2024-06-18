import { RequestHandler } from "express";
import { WebsocketRequestHandler } from "express-ws";
import TermService from "../service/term-service";
import { UserServiceError } from "../types/error";
import { ResizeTermScheme } from "../types/request";

const create: RequestHandler = (req, res) => {
  const user = req.session.user!;
  if (user.session) {
    res.status(409)
    res.send({ error: "Session exists" });
  }
  try {
    const uuid = TermService.create(req.session.user.username);
    user.session = uuid;
    res.status(201);
    res.end();
  } catch (error) {
    if (error instanceof UserServiceError) {
      res.status(error.code);
      res.send({ error: error.message });
      return;
    } else {
      res.status(500);
      res.send({ error: "Internal server error" });
      return;
    }
  }
}

const resize: RequestHandler = (req, res) => {
  const user = req.session.user!;
  let schema = ResizeTermScheme.safeParse(req.body);
  if (!schema.success) {
    res.status(400)
    res.send({ error: "Invalid input" });
    return;
  }
  let body = schema.data;

  let uuid = user.session;
  if (uuid === undefined || typeof uuid !== 'string') {
    res.status(400)
    res.send({ error: "Session invalid or not existent" });
    return;
  }

  TermService.resize(uuid, { cols: body.cols!, rows: body.rows! });
  res.status(200);
  res.end();
}

const ws: WebsocketRequestHandler = (ws, req) => {
  const user = req.session.user!;
  let uuid = user.session;
  if (uuid === undefined || typeof uuid !== 'string') {
    console.log("[ERROR] Session invalid or not existent");
    ws.close();
    return;
  }
  try {
    TermService.ws(uuid, ws, () => {
      user.session = undefined;
    })
  } catch (error) {
    console.log("[ERROR]", error);
    ws.close();
  }
}

const TermController = {
  create,
  resize,
  ws
}

export default TermController;

