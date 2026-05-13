import { redirect } from 'react-router'

import { UnauthenticatedError } from '@/lib/api/errors'

import { buildUrl } from '@/router'

export const loginRequired = <Args extends unknown[], ReturnType>(
  fn: (...args: Args) => Promise<ReturnType>
) => {
  return async (...args: Args) => {
    try {
      return await fn(...args)
    } catch (error) {
      if (error instanceof UnauthenticatedError) {
        throw redirect(buildUrl('/login'))
      } else {
        throw error
      }
    }
  }
}
