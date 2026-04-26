import { createContext } from 'react'

import { CloseNotification, ShowNotification } from '@/lib/models/Notification'

const NotificationsContext = createContext<{
  show: ShowNotification
  close: CloseNotification
}>(null!)

export default NotificationsContext
