import { z } from "zod";

export async function sleep(ms: number) {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve(null);
    }, ms);
  })
}

export function randomString(n: number): string {
  const charSet = 'abcdefghijklmnopqrstuvwxyz012345';
  var str = '';
  for (let i = 0; i < n; i++) {
    const j = Math.floor(Math.random() * charSet.length);
    str += charSet[j];
  }
  return str;
}

export const UserSessionResponseSchema = z.array(z.string())
export type UserSesssionResponse = z.infer<typeof UserSessionResponseSchema>;

