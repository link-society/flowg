import {
  ShowNotificationOptions,
  useNotifications,
} from '@toolpad/core/useNotifications'

export const useNotify = () => {
  const notifications = useNotifications()

  type Severity = ShowNotificationOptions['severity']
  const notify = (severity: Severity, message: string) => {
    notifications.show(message, {
      severity,
      autoHideDuration: 3000,
    })
  }

  return {
    info: (message: string) => notify('info', message),
    warning: (message: string) => notify('warning', message),
    error: (message: string) => notify('error', message),
    success: (message: string) => notify('success', message),
  }
}
