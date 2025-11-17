import Button from '@mui/material/Button'
import { useDialogs } from '@toolpad/core/useDialogs'

import AddIcon from '@mui/icons-material/Add'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import DialogNewForwarder from '@/components/DialogNewForwarder'

type ButtonNewForwarderProps = Readonly<{
  onForwarderCreated: (name: string) => void
}>

const ButtonNewForwarder = (props: ButtonNewForwarderProps) => {
  const dialogs = useDialogs()
  const notify = useNotify()

  const [handleClick] = useApiOperation(async () => {
    const forwarderName = (await dialogs.open(DialogNewForwarder)) as
      | string
      | null
    if (forwarderName !== null) {
      props.onForwarderCreated(forwarderName)

      notify.success('Forwarder created')
    }
  }, [props.onForwarderCreated])

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

export default ButtonNewForwarder
