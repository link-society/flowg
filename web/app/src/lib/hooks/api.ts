import { useCallback, useState } from 'react'
import { useNavigate } from 'react-router'

import { PermissionDeniedError, UnauthenticatedError } from '@/lib/api/errors'

import { useNotify } from '@/lib/hooks/notify'

export function useApiOperation<Args extends unknown[]>(
  fn: (...args: Args) => Promise<void>,
  deps: unknown[]
): [(...args: Args) => Promise<void>, boolean] {
  const [loading, setLoading] = useState(false)

  const navigate = useNavigate()
  const notify = useNotify()

  const cb = useCallback(
    async (...args: Args) => {
      setLoading(true)

      try {
        await fn(...args)
      } catch (error) {
        console.error(error)

        if (error instanceof UnauthenticatedError) {
          notify.error('Session expired')
          navigate('/web/login')
        } else if (error instanceof PermissionDeniedError) {
          notify.error('Permission denied')
        } else {
          notify.error('Unknown error')
        }
      }

      setLoading(false)
    },
    [fn, navigate, notify, ...deps]
  )

  return [cb, loading]
}
