import { useDialogs } from '@toolpad/core/useDialogs'
import { useNotifications } from '@toolpad/core/useNotifications'
import { useConfig } from '@/lib/context/config'
import { useApiOperation } from '@/lib/hooks/api'

import AddIcon from '@mui/icons-material/Add'

import Button from '@mui/material/Button'

import { NewTransformerModal } from './modal'

type NewTransformerButtonProps = {
  onTransformerCreated: (name: string) => void
}

export const NewTransformerButton = (props: NewTransformerButtonProps) => {
  const dialogs = useDialogs()
  const notifications = useNotifications()
  const config = useConfig()

  const [handleClick] = useApiOperation(
    async () => {
      const transformerName = await dialogs.open(NewTransformerModal) as string | null
      if (transformerName !== null) {
        props.onTransformerCreated(transformerName)

        notifications.show('Transformer created', {
          severity: 'success',
          autoHideDuration: config.notifications?.autoHideDuration,
        })
      }
    },
    [props.onTransformerCreated],
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
