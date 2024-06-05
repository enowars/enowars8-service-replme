import { z } from "zod";

export const CreateUserScheme = z.object({
  username: z.string().min(4).max(64).regex(/^[a-zA-Z0-5]*$/),
  password: z.string().min(4).max(64).regex(/^[a-zA-Z0-9]*$/)
});

export type CreateUser = z.infer<typeof CreateUserScheme>;

export const LoginUserScheme = z.object({
  username: z.string().min(4).max(64).regex(/^[a-zA-Z0-5]*$/),
  password: z.string().min(4).max(64).regex(/^[a-zA-Z0-9]*$/)
});

export type LoginUser = z.infer<typeof LoginUserScheme>;

export const ResizeTermScheme = z.object({
  cols: z.number().positive(),
  rows: z.number().positive()
});

export type ResizeTerm = z.infer<typeof ResizeTermScheme>;
