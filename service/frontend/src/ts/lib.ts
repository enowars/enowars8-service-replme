
import * as z from "zod";

const LoginsSchema = z.record(z.string(), z.string());
type Logins = z.infer<typeof LoginsSchema>;

export function getLogins(): Logins {
  const logins = LoginsSchema.safeParse(JSON.parse(localStorage.getItem('logins') ?? "{}"));
  return logins.success ? logins.data : {};
}

export function getPassword(username: string): string | undefined {
  return getLogins()[username]
}

export function addLogin(username: string, password: string) {
  const logins = getLogins();
  logins[username] = password;
  localStorage.setItem('logins', JSON.stringify(logins));
}

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
