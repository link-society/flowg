import { useDialogs } from '@toolpad/core/useDialogs'
import { useNavigate } from 'react-router-dom'
import { useNotifications } from '@toolpad/core/useNotifications'
import { useConfig } from '@/lib/context/config'

import AddIcon from '@mui/icons-material/Add'

import Button from '@mui/material/Button'

import { NewTransformerModal } from './modal'

import { UnauthenticatedError, PermissionDeniedError } from '@/lib/api/errors'

export const NewTransformerButton = () => {
  const dialogs = useDialogs()
  const navigate = useNavigate()
  const notifications = useNotifications()
  const config = useConfig()

  const handleClick = async () => {
    try {
      const transformerName = await dialogs.open(NewTransformerModal) as string | null
      if (transformerName !== null) {
        window.location.pathname = `/web/transformers/${transformerName}`

        notifications.show('Transformer created', {
          severity: 'success',
          autoHideDuration: config.notifications?.autoHideDuration,
        })
      }
    }
    catch (error) {
      if (error instanceof UnauthenticatedError) {
        notifications.show('Session expired', {
          severity: 'error',
          autoHideDuration: config.notifications?.autoHideDuration,
        })
        navigate('/web/login')
      }
      else if (error instanceof PermissionDeniedError) {
        notifications.show('Permission denied', {
          severity: 'error',
          autoHideDuration: config.notifications?.autoHideDuration,
        })
      }
      else {
        notifications.show('Unknown error', {
          severity: 'error',
          autoHideDuration: config.notifications?.autoHideDuration,
        })
      }

      console.error(error)
    }
  }

  return (
    <Button
      variant="contained"
      color="primary"
      size="small"
      onClick={() => handleClick()}
      startIcon={<AddIcon />}
    >
      New
    </Button>
  )
}
