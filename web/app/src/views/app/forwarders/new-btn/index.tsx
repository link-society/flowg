import { useDialogs } from '@toolpad/core/useDialogs'
import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import AddIcon from '@mui/icons-material/Add'

import Button from '@mui/material/Button'

import { NewForwarderModal } from './modal'

type NewForwarderButtonProps = Readonly<{
  onForwarderCreated: (name: string) => void
}>

export const NewForwarderButton = (props: NewForwarderButtonProps) => {
  const dialogs = useDialogs()
  const notify = useNotify()

  const [handleClick] = useApiOperation(
    async () => {
      const forwarderName = await dialogs.open(NewForwarderModal) as string | null
      if (forwarderName !== null) {
        props.onForwarderCreated(forwarderName)

        notify.success('Forwarder created')
      }
    },
    [props.onForwarderCreated],
  )

  return (
    <Button
      id="btn:forwarders.create"
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
