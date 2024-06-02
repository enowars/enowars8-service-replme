import { z } from "zod";

export const LoginUserScheme = z.object({
  username: z.string().min(4).max(64).regex(/^[a-zA-Z0-5]*$/),
  password: z.string().min(4).max(64).regex(/^[a-zA-Z0-9]*$/)
});

export type LoginUser = z.infer<typeof LoginUserScheme>;

export const CreateUserScheme = z.object({
  username: z.string().min(4).max(64).regex(/^[a-zA-Z0-5]*$/),
  password: z.string().min(4).max(64).regex(/^[a-zA-Z0-9]*$/)
});

export type CreateUser = z.infer<typeof CreateUserScheme>;
