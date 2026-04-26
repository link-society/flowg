import { useContext } from 'react'

import NotificationsContext from '@/lib/context/notifications'

export const useNotifications = () => useContext(NotificationsContext)
