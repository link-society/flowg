import { useCallback, useMemo, useState } from 'react'

import NotificationsContext from '@/lib/context/notifications'

import {
  CloseNotification,
  NotificationsState,
  ShowNotification,
} from '@/lib/models/Notification'

import Notifications from '@/components/Notifications/component'

import { NotificationsProviderProps, RootPropsContext } from './context'

let nextId = 0
const generateId = () => {
  const id = nextId
  nextId += 1
  return id
}

const NotificationsProvider = (props: NotificationsProviderProps) => {
  const { children } = props
  const [state, setState] = useState<NotificationsState>({ queue: [] })

  const show = useCallback<ShowNotification>((message, options = {}) => {
    const notificationKey =
      options.key ?? `::flowg-internal::notification::${generateId()}`
    setState((prev) => {
      if (
        prev.queue.some((entry) => entry.notificationKey === notificationKey)
      ) {
        return prev
      }

      return {
        ...prev,
        queue: [
          ...prev.queue,
          {
            message,
            options,
            notificationKey,
            open: true,
          },
        ],
      }
    })

    return notificationKey
  }, [])

  const close = useCallback<CloseNotification>((key) => {
    setState((prev) => ({
      ...prev,
      queue: prev.queue.filter((n) => n.notificationKey !== key),
    }))
  }, [])

  const contextValue = useMemo(() => ({ show, close }), [show, close])

  return (
    <RootPropsContext.Provider value={props}>
      <NotificationsContext.Provider value={contextValue}>
        {children}
        <Notifications state={state} />
      </NotificationsContext.Provider>
    </RootPropsContext.Provider>
  )
}

export default NotificationsProvider
