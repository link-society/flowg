import { useDialogs } from '@toolpad/core/useDialogs'
import { useNotifications } from '@toolpad/core/useNotifications'
import { useConfig } from '@/lib/context/config'
import { useApiOperation } from '@/lib/hooks/api'

import AddIcon from '@mui/icons-material/Add'

import Button from '@mui/material/Button'

import { NewPipelineModal } from './modal'

type NewPipelineButtonProps = {
  onPipelineCreated: (name: string) => void
}

export const NewPipelineButton = (props: NewPipelineButtonProps) => {
  const dialogs = useDialogs()
  const notifications = useNotifications()
  const config = useConfig()

  const [handleClick] = useApiOperation(
    async () => {
      const pipelineName = await dialogs.open(NewPipelineModal) as string | null
      if (pipelineName !== null) {
        props.onPipelineCreated(pipelineName)

        notifications.show('Pipeline created', {
          severity: 'success',
          autoHideDuration: config.notifications?.autoHideDuration,
        })
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
