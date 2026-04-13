import Notification from '@/components/Notification/component'

import { NotificationsProps } from './types'

const Notifications = ({ state }: NotificationsProps) => {
  const currentNotification = state.queue[0] ?? null

  return currentNotification ? (
    <Notification
      {...currentNotification}
      badge={state.queue.length > 1 ? String(state.queue.length) : null}
    />
  ) : null
}

export default Notifications
