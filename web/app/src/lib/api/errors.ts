export class ApiError extends Error {}

export class InvalidResponseError extends ApiError {
  body: string
  statusCode: number
  ok: boolean

  constructor(message: string, body: string, statusCode: number, ok: boolean) {
    super(message)

    this.body = body
    this.statusCode = statusCode
    this.ok = ok
  }
}

export class UnauthenticatedError extends ApiError {}

export class PermissionDeniedError extends ApiError {}
