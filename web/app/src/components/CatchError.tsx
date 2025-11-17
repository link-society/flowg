import { ReactNode } from 'react'
import { Navigate, useAsyncError } from 'react-router'

type CatchErrorProps<E extends Error> = {
  errorType: new (...args: any[]) => E
  fallback: ReactNode
}

const CatchError = <E extends Error>(props: CatchErrorProps<E>) => {
  const error = useAsyncError()

  if (error instanceof props.errorType) {
    return <Navigate to="/web/login" />
  } else if (props.fallback) {
    return <>{props.fallback}</>
  } else {
    throw error
  }
}

export default CatchError
