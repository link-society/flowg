import { Navigate, useAsyncError } from 'react-router'

import { buildUrl } from '@/router'

import { CatchErrorProps } from './types'

const CatchError = <E extends Error>(props: CatchErrorProps<E>) => {
  const error = useAsyncError()

  if (error instanceof props.errorType) {
    return <Navigate to={buildUrl('/login')} />
  } else if (props.fallback) {
    return <>{props.fallback}</>
  } else {
    throw error
  }
}

export default CatchError
