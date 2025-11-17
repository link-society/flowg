import { ReactNode } from 'react'
import { Navigate, useAsyncError } from 'react-router'

type CatchErrorProps<E extends Error> = {
  errorType: new (...args: any[]) => E
  fallback: ReactNode
}

function CatchError<E extends Error>({
  errorType,
  fallback,
}: CatchErrorProps<E>) {
  const error = useAsyncError()

  if (error instanceof errorType) {
    return <Navigate to="/web/login" />
  } else if (fallback) {
    return <>{fallback}</>
  } else {
    throw error
  }
}

export default CatchError
