import { createContext } from 'react'

import { SnackbarProps } from '@mui/material/Snackbar'

type NotificationsProviderSlotProps = {
  snackbar: SnackbarProps
}

type NotificationsProviderSlots = {
  snackbar: React.ElementType
}

export type NotificationsProviderProps = {
  children?: React.ReactNode
  slots?: Partial<NotificationsProviderSlots>
  slotProps?: Partial<NotificationsProviderSlotProps>
}

export const RootPropsContext = createContext<NotificationsProviderProps>(null!)
