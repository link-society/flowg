export type ApiErrorResponse = {
  status: string
  error: string
  code?: number
  context?: Record<string, string>
}
