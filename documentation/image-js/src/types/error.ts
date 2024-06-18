
export class UserServiceError extends Error {
  code: number
  constructor(code: number, message: string) {
    super(message);
    this.code = code;
  }
}

export class TermServiceError extends Error {
  code: number
  constructor(code: number, message: string) {
    super(message);
    this.code = code;
  }
};
