import { useDialogs } from '@toolpad/core/useDialogs'
import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import AddIcon from '@mui/icons-material/Add'

import Button from '@mui/material/Button'

import { NewStreamModal } from './modal'

type NewStreamButtonProps = Readonly<{
  onStreamCreated: (name: string) => void
}>

export const NewStreamButton = (props: NewStreamButtonProps) => {
  const dialogs = useDialogs()
  const notify = useNotify()

  const [handleClick] = useApiOperation(
    async () => {
      const streamName = await dialogs.open(NewStreamModal) as string | null
      if (streamName !== null) {
        props.onStreamCreated(streamName)

        notify.success('Stream created')
      }
    },
    [props.onStreamCreated],
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