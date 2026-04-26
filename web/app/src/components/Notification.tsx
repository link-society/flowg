import { Fragment, useCallback, useContext } from 'react'

import useSlotProps from '@mui/utils/useSlotProps'

import Alert from '@mui/material/Alert'
import Badge from '@mui/material/Badge'
import Button from '@mui/material/Button'
import IconButton from '@mui/material/IconButton'
import Snackbar, { SnackbarCloseReason } from '@mui/material/Snackbar'
import SnackbarContent from '@mui/material/SnackbarContent'
import { CloseReason } from '@mui/material/SpeedDial'

import CloseIcon from '@mui/icons-material/Close'

import NotificationsContext from '@/lib/context/notifications'

import { ShowNotificationOptions } from '@/lib/models/Notification'

import { RootPropsContext } from './NotificationsProvider'

type NotificationProps = {
  notificationKey: string
  badge: string | null
  open: boolean
  message: React.ReactNode
  options: ShowNotificationOptions
}

const Notification = ({
  notificationKey,
  badge,
  open,
  message,
  options,
}: NotificationProps) => {
  const { close } = useContext(NotificationsContext)
  const { severity, actionText, onAction, autoHideDuration } = options

  const handleClose = useCallback(
    (_event: unknown, reason?: CloseReason | SnackbarCloseReason) => {
      if (reason === 'clickaway') {
        return
      }

      close(notificationKey)
    },
    [close, notificationKey]
  )

  const action = (
    <Fragment>
      {onAction ? (
        <Button color="inherit" size="small" onClick={onAction}>
          {actionText ?? 'Action'}
        </Button>
      ) : null}
      <IconButton
        size="small"
        aria-label="close"
        title="Close"
        color="inherit"
        onClick={handleClose}
      >
        <CloseIcon fontSize="small" />
      </IconButton>
    </Fragment>
  )

  const props = useContext(RootPropsContext)
  const SnackbarComponent = props?.slots?.snackbar ?? Snackbar
  const snackbarSlotProps = useSlotProps({
    elementType: SnackbarComponent,
    ownerState: props,
    externalSlotProps: props?.slotProps?.snackbar,
    additionalProps: {
      open,
      autoHideDuration,
      onClose: handleClose,
      action,
    },
  })

  return (
    <SnackbarComponent key={notificationKey} {...snackbarSlotProps}>
      <Badge badgeContent={badge} color="primary" sx={{ width: '100%' }}>
        {severity ? (
          <Alert severity={severity} sx={{ width: '100%' }} action={action}>
            {message}
          </Alert>
        ) : (
          <SnackbarContent message={message} action={action} />
        )}
      </Badge>
    </SnackbarComponent>
  )
}

export default Notification
