import { redirect } from 'react-router-dom'

import { UnauthenticatedError } from '@/lib/api/errors'

export const loginRequired = <Args extends unknown[], ReturnType>(
  fn: (...args: Args) => Promise<ReturnType>
) => {
  return async (...args: Args) => {
    try {
      return await fn(...args)
    }
    catch (error) {
      if (error instanceof UnauthenticatedError) {
        throw redirect('/web/login')
      }
      else {
        throw error
      }
    }
  }
}
