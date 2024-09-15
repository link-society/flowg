import { Await, AwaitProps, Navigate, useAsyncError } from 'react-router-dom'

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

export const AuthenticatedAwait = (props: AwaitProps) => (
  <Await
    resolve={props.resolve}
    errorElement={<CatchUnauthenticatedError fallback={props.errorElement} />}
  >
    {props.children}
  </Await>
)
