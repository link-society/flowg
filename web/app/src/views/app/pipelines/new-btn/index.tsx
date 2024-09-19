import { useCallback } from 'react'
import { useDialogs } from '@toolpad/core/useDialogs'
import { useNavigate } from 'react-router-dom'
import { useNotifications } from '@toolpad/core/useNotifications'
import { useConfig } from '@/lib/context/config'

import AddIcon from '@mui/icons-material/Add'

import Button from '@mui/material/Button'

import { NewPipelineModal } from './modal'

import { UnauthenticatedError, PermissionDeniedError } from '@/lib/api/errors'

type NewPipelineButtonProps = {
  onPipelineCreated: (name: string) => void
}

export const NewPipelineButton = (props: NewPipelineButtonProps) => {
  const dialogs = useDialogs()
  const navigate = useNavigate()
  const notifications = useNotifications()
  const config = useConfig()

  const handleClick = useCallback(
    async () => {
      try {
        const pipelineName = await dialogs.open(NewPipelineModal) as string | null
        if (pipelineName !== null) {
          props.onPipelineCreated(pipelineName)

          notifications.show('Pipeline created', {
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
    },
    [props.onPipelineCreated],
  )

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
