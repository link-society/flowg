import { ApiErrorResponse } from '@/lib/api/response'

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

export class ApplicationError extends ApiError {
  appcode: number
  context: Record<string, string>

  constructor(resp: ApiErrorResponse, prefix: string = '') {
    let message = resp.error
    if (message.startsWith(prefix)) {
      message = message.slice(prefix.length).trim()
    }

    super(message)

    this.appcode = resp.code ?? 0
    this.context = resp.context ?? {}
  }
}

export class UnauthenticatedError extends ApplicationError {
  constructor(resp: ApiErrorResponse) {
    super(resp, 'unauthenticated: ')
  }
}

export class PermissionDeniedError extends ApplicationError {
  constructor(resp: ApiErrorResponse) {
    super(resp, 'permission denied: ')
  }
}

export class InvalidArgumentError extends ApplicationError {
  constructor(resp: ApiErrorResponse) {
    super(resp, 'invalid argument: ')
  }
}
