import { Await, AwaitProps } from 'react-router'

import { UnauthenticatedError } from '@/lib/api/errors'

import CatchError from '@/components/CatchError'

function AuthenticatedAwait<T>(props: AwaitProps<T>) {
  return (
    <Await
      resolve={props.resolve}
      errorElement={
        <CatchError
          errorType={UnauthenticatedError}
          fallback={props.errorElement}
        />
      }
    >
      {props.children}
    </Await>
  )
}

export default AuthenticatedAwait
