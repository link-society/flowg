import { Await, AwaitProps, Navigate, useAsyncError } from 'react-router'

import { UnauthenticatedError } from '@/lib/api/errors'

const CatchUnauthenticatedError = (props: { fallback: React.ReactNode }) => {
  const error = useAsyncError()

  if (error instanceof UnauthenticatedError) {
    return <Navigate to="/web/login" />
  } else if (props.fallback) {
    return <>{props.fallback}</>
  } else {
    throw error
  }
}

export function AuthenticatedAwait<T>(props: AwaitProps<T>) {
  return (
    <Await
      resolve={props.resolve}
      errorElement={<CatchUnauthenticatedError fallback={props.errorElement} />}
    >
      {props.children}
    </Await>
  )
}
