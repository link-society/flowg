import { useCallback, useState } from 'react'
import { useNavigate } from 'react-router'
import { useNotify } from '@/lib/hooks/notify'

import { UnauthenticatedError, PermissionDeniedError } from '@/lib/api/errors'

export function useApiOperation<Args extends unknown[]>(
  fn: (...args: Args) => Promise<void>,
  deps: unknown[],
): [(...args: Args) => Promise<void>, boolean] {
  const [loading, setLoading] = useState(false)

  const navigate = useNavigate()
  const notify = useNotify()

  const cb = useCallback(
    async (...args: Args) => {
      setLoading(true)

      try {
        await fn(...args)
      }
      catch (error) {
        if (error instanceof UnauthenticatedError) {
          notify.error('Session expired')
          navigate('/web/login')
        }
        else if (error instanceof PermissionDeniedError) {
          notify.error('Permission denied')
        }
        else {
          notify.error('Unknown error')
        }

        console.error(error)
      }

      setLoading(false)
    },
    [fn, ...deps],
  )

  return [cb, loading]
}
