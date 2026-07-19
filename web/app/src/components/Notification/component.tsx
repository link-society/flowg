import { Fragment, useCallback, useContext } from 'react'
import { useTranslation } from 'react-i18next'

import useSlotProps from '@mui/utils/useSlotProps'

import Button from '@mui/material/Button'
import IconButton from '@mui/material/IconButton'
import Snackbar, { SnackbarCloseReason } from '@mui/material/Snackbar'
import SnackbarContent from '@mui/material/SnackbarContent'
import { CloseReason } from '@mui/material/SpeedDial'

import CloseIcon from '@mui/icons-material/Close'

import NotificationsContext from '@/lib/context/notifications'

import { RootPropsContext } from '@/components/NotificationsProvider/context'

import { NotificationAlert, NotificationBadge } from './styles'
import { NotificationProps } from './types'

const Notification = ({
  notificationKey,
  badge,
  open,
  message,
  options,
}: NotificationProps) => {
  const { t } = useTranslation()
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
          {actionText ?? t('components.notification.action')}
        </Button>
      ) : null}
      <IconButton
        size="small"
        aria-label={t('common.actions.close')}
        title={t('common.actions.close')}
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
      <NotificationBadge badgeContent={badge} color="primary">
        {severity ? (
          <NotificationAlert severity={severity} action={action}>
            {message}
          </NotificationAlert>
        ) : (
          <SnackbarContent message={message} action={action} />
        )}
      </NotificationBadge>
    </SnackbarComponent>
  )
}

export default Notification
