import { useDialogs } from '@toolpad/core/useDialogs'
import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import AddIcon from '@mui/icons-material/Add'

import Button from '@mui/material/Button'

import { NewTransformerModal } from './modal'

type NewTransformerButtonProps = Readonly<{
  onTransformerCreated: (name: string) => void
}>

export const NewTransformerButton = (props: NewTransformerButtonProps) => {
  const dialogs = useDialogs()
  const notify = useNotify()

  const [handleClick] = useApiOperation(
    async () => {
      const transformerName = await dialogs.open(NewTransformerModal) as string | null
      if (transformerName !== null) {
        props.onTransformerCreated(transformerName)

        notify.success('Transformer created')
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
