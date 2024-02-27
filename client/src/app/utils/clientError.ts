export class ClientError extends Error {
  code: number;
  cause?: string;
  constructor(message: string, code: number = 500, cause?: string) {
    super(message);
    this.code = code;
    this.cause = cause;
  }
}
