import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function randomString(n: number): string {
  const charSet = "abcdefghijklmnopqrstuvwxyz012345";
  var str = "";
  for (let i = 0; i < n; i++) {
    const j = Math.floor(Math.random() * charSet.length);
    str += charSet[j];
  }
  return str;
}

export async function sleep(ms: number) {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve(null);
    }, ms);
  });
}
