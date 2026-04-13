import { ShowNotificationOptions } from '@/lib/models/Notification'

export type NotificationProps = {
  notificationKey: string
  badge: string | null
  open: boolean
  message: React.ReactNode
  options: ShowNotificationOptions
}
