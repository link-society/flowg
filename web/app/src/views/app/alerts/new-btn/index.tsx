import { useDialogs } from '@toolpad/core/useDialogs'
import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import AddIcon from '@mui/icons-material/Add'

import Button from '@mui/material/Button'

import { NewAlertModal } from './modal'

type NewAlertButtonProps = Readonly<{
  onAlertCreated: (name: string) => void
}>

export const NewAlertButton = (props: NewAlertButtonProps) => {
  const dialogs = useDialogs()
  const notify = useNotify()

  const [handleClick] = useApiOperation(
    async () => {
      const alertName = await dialogs.open(NewAlertModal) as string | null
      if (alertName !== null) {
        props.onAlertCreated(alertName)

        notify.success('Alert created')
      }
    },
    [props.onAlertCreated],
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
