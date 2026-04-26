export interface ShowNotificationOptions {
  key?: string
  severity?: 'info' | 'warning' | 'error' | 'success'
  autoHideDuration?: number
  actionText?: React.ReactNode
  onAction?: () => void
}

export interface ShowNotification {
  (message: React.ReactNode, options?: ShowNotificationOptions): string
}

export interface CloseNotification {
  (key: string): void
}

export interface NotificationQueueEntry {
  notificationKey: string
  options: ShowNotificationOptions
  open: boolean
  message: React.ReactNode
}

export interface NotificationsState {
  queue: NotificationQueueEntry[]
}
