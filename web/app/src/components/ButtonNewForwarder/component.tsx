import Button from '@mui/material/Button'

import AddIcon from '@mui/icons-material/Add'

import { useApiOperation } from '@/lib/hooks/api'
import { useDialogs } from '@/lib/hooks/dialogs'
import { useNotify } from '@/lib/hooks/notify'

import DialogNewForwarder from '@/components/DialogNewForwarder/component'

import { ButtonNewForwarderProps } from './types'

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
