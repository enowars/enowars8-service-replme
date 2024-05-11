import fs from "node:fs";
import child_process from "node:child_process";

import { UserServiceError } from "types/error";

interface UserEntry {
  hash: string;
  type: string;
  salt: string;
  pw: string;
}

function findUserHash(username: string): string | undefined {
  const shadow = fs.readFileSync('/etc/shadow', 'utf8').split('\n');

  const hash = shadow.find((user) => {
    return user.startsWith(username)
  })

  return hash;
}

function extractUserEntry(line: string): UserEntry {
  const [, hash] = line.split(':')
  const _hash = hash.split('$')
  if (_hash.length < 4) {
    throw new UserServiceError(403, "Forbidden");
  }
  let type: string, salt: string, pw: string = '';
  if (hash[1] === 'y') {
    throw new UserServiceError(403, "Forbidden");
  } else {
    [, type, salt, pw] = _hash;
  }
  return {
    hash,
    type,
    salt,
    pw
  }
}

function validateUserPassword(entry: UserEntry, password: string): boolean {
  const buf = child_process.execSync(
    `openssl passwd -${entry.type} -salt ${entry.salt} ${password}`
  );
  const result = buf.toString('utf8').trim();
  return result === entry.hash;
}

function login(username: string, password: string) {
  const hash = findUserHash(username);

  if (!hash) {
    throw new UserServiceError(404, "Not found");
  }

  const user = extractUserEntry(hash);
  return validateUserPassword(user, password);
}

function createUserEntry(username: string, password: string) {
  try {
    child_process.execSync(`adduser -D ${username} -s /bin/zsh`);
    child_process.execSync(`echo "${username}:${password}" | chpasswd`);
  } catch (error) {
    throw new UserServiceError(500, "Internal server error");
  }
}

function create(username: string, password: string): "ok" | "created" {

  const hash = findUserHash(username);

  if (hash) {
    const user = extractUserEntry(hash);
    const valid =  validateUserPassword(user, password);
    if (!valid) {
      throw new UserServiceError(401, "Invalid credentials");
    }
    return "ok";
  }

  createUserEntry(username, password);

  return "created";
}

const UserService = {
  login,
  create
}

export default UserService;

