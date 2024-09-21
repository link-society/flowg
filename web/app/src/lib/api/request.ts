import { EventSourcePlus } from 'event-source-plus'

import * as errors from '@/lib/api/errors'

type ApiMethod = 'GET' | 'POST' | 'PUT' | 'DELETE'

type ApiErrorResponse = {
  status: string
  error: string
}

type ApiRequestResult<R> = {
  body: R
  response: Response
}

type RequestOptions<B> = {
  path: string
  searchParams?: URLSearchParams,
  body?: B
}

const request = async<B, R extends { success: boolean }>(
  method: ApiMethod,
  { path, searchParams, body }: RequestOptions<B>,
): Promise<ApiRequestResult<R>> => {
  let authHeader = ''

  const token = localStorage.getItem('token')
  if (token !== null) {
    authHeader = `Bearer jwt:${token}`
  }

  const response = await fetch(
    `${path}?${searchParams?.toString() ?? ''}`,
    {
      method,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': authHeader,
      },
      body: JSON.stringify(body),
    },
  )
  const responseBody = await response.json()

  if (!response.ok) {
    const content = responseBody as ApiErrorResponse

    switch (content.status) {
      case 'UNAUTHENTICATED':
        throw new errors.UnauthenticatedError(content.error)

      case 'PERMISSION_DENIED':
        throw new errors.PermissionDeniedError(content.error)

      default:
        throw new errors.ApiError(content.error)
    }
  }

  if (!responseBody.success) {
    throw new errors.ApiError('Request failed')
  }

  return {
    body: responseBody as R,
    response,
  }
}

export const GET = async<R extends { success: boolean }>(
  options: RequestOptions<never>,
): Promise<ApiRequestResult<R>> => {
  return request('GET', options)
}

export const POST = async<B, R extends { success: boolean }>(
  options: RequestOptions<B>,
): Promise<ApiRequestResult<R>> => {
  return request('POST', options)
}

export const PUT = async<B, R extends { success: boolean }>(
  options: RequestOptions<B>,
): Promise<ApiRequestResult<R>> => {
  return request('PUT', options)
}

export const DELETE = async<R extends { success: boolean }>(
  options: RequestOptions<never>,
): Promise<ApiRequestResult<R>> => {
  return request('DELETE', options)
}

export const SSE = (
  options: RequestOptions<never>,
) => {
  const eventSource = new EventSourcePlus(
    `${options.path}?${options.searchParams?.toString() ?? ''}`,
    {
      maxRetryCount: 0,
      headers: {
        'Authorization': `Bearer jwt:${localStorage.getItem('token')}`,
      },
    },
  )

  const messages = new EventTarget()
  const control = new EventTarget()

  const cancelScope = eventSource.listen({
    onRequestError({ error }) {
      cancelScope.abort()
      control.dispatchEvent(new CustomEvent('error', { detail: error }))
    },
    async onResponseError({ response }) {
      try {
        const responseBody = await response.json()
        const content = responseBody as ApiErrorResponse

        switch (content.status) {
          case 'UNAUTHENTICATED':
            throw new errors.UnauthenticatedError(content.error)

          case 'PERMISSION_DENIED':
            throw new errors.PermissionDeniedError(content.error)

          default:
            throw new errors.ApiError(content.error)
        }
      }
      catch (error) {
        cancelScope.abort()
        control.dispatchEvent(new CustomEvent('error', { detail: error }))
      }
    },
    onMessage(msg) {
      const detail = {
        id: msg.id,
        data: msg.data,
        retry: msg.retry,
      }

      messages.dispatchEvent(new CustomEvent(msg.event, { detail }))
    }
  })

  return {
    messages,
    control,
    close() {
      cancelScope.abort()
    }
  }
}
