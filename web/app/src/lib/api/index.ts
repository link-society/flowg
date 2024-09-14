type ApiMethod = 'GET' | 'POST' | 'PUT' | 'DELETE'

type ApiErrorResponse = {
  status: string
  error: string
}

type ApiRequestResult<R> = {
  body: R
  response: Response
}

export class ApiError extends Error {}
export class UnauthenticatedError extends ApiError {}

const request = async<B, R>(
  method: ApiMethod,
  path: string,
  body?: B,
): Promise<ApiRequestResult<R>> => {
  let authHeader = ''

  const token = localStorage.getItem('token')
  if (token !== null) {
    authHeader = `Bearer ${token}`
  }

  const response = await fetch(path, {
    method,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': authHeader,
    },
    body: JSON.stringify(body),
  })
  const responseBody = await response.json()

  if (!response.ok) {
    const content = responseBody as ApiErrorResponse

    switch (content.status) {
      case 'UNAUTHENTICATED':
        throw new UnauthenticatedError(content.error)

      default:
        throw new ApiError(content.error)
    }
  }

  return {
    body: responseBody as R,
    response,
  }
}

export const GET = async<R>(path: string): Promise<ApiRequestResult<R>> => {
  return request('GET', path)
}

export const POST = async<B, R>(path: string, body: B): Promise<ApiRequestResult<R>> => {
  return request('POST', path, body)
}

export const PUT = async<B, R>(path: string, body: B): Promise<ApiRequestResult<R>> => {
  return request('PUT', path, body)
}

export const DELETE = async<R>(path: string): Promise<ApiRequestResult<R>> => {
  return request('DELETE', path)
}
