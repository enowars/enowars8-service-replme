import { RequestHandler } from "express";
import UserService from "../service/user-service";
import { UserServiceError } from "../types/error";
import { CreateUserScheme, LoginUserScheme } from "../types/request";

const register: RequestHandler = (req, res) => {
  let schema = CreateUserScheme.safeParse(req.body);
  if (!schema.success) {
    res.status(400)
    res.send({ error: "Invalid input" });
    return;
  }
  let body = schema.data;

  try {
    if (UserService.create(body.username, body.password) === "ok") {
      req.session.user = {
        username: body.username,
      }
      res.status(200);
      res.send({ success: "Ok" });
      return;
    } else {
      req.session.user = {
        username: body.username,
      }
      res.status(201);
      res.send({ success: "Created" });
      return;
    }
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

const login: RequestHandler = (req, res) => {
  let schema = LoginUserScheme.safeParse(req.body);
  if (!schema.success) {
    res.status(400)
    res.send({ error: "Invalid input" });
    return;
  }
  let body = schema.data;

  try {
    if (UserService.login(body.username, body.password)) {
      req.session.user = {
        username: body.username,
      }
      res.status(200);
      res.send({ success: "Ok" });
      return;
    } else {
      res.status(401);
      res.send({ error: "Invalid credentials" });
      return;
    }
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

const AuthController = {
  register,
  login,
}

export default AuthController;

