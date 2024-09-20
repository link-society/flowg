import { useCallback, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useNotifications } from '@toolpad/core/useNotifications'
import { useConfig } from '@/lib/context/config'

import { UnauthenticatedError, PermissionDeniedError } from '@/lib/api/errors'

export function useApiOperation<Args extends unknown[]>(
  fn: (...args: Args) => Promise<void>,
  deps: unknown[],
): [(...args: Args) => Promise<void>, boolean] {
  const [loading, setLoading] = useState(false)

  const navigate = useNavigate()
  const notifications = useNotifications()
  const config = useConfig()

  const cb = useCallback(
    async (...args: Args) => {
      setLoading(true)

      try {
        await fn(...args)
      }
      catch (error) {
        if (error instanceof UnauthenticatedError) {
          notifications.show('Session expired', {
            severity: 'error',
            autoHideDuration: config.notifications?.autoHideDuration,
          })
          navigate('/web/login')
        }
        else if (error instanceof PermissionDeniedError) {
          notifications.show('Permission denied', {
            severity: 'error',
            autoHideDuration: config.notifications?.autoHideDuration,
          })
        }
        else {
          notifications.show('Unknown error', {
            severity: 'error',
            autoHideDuration: config.notifications?.autoHideDuration,
          })
        }

        console.error(error)
      }

      setLoading(false)
    },
    [fn, ...deps],
  )

  return [cb, loading]
}
