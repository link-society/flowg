import Button from '@mui/material/Button'
import { useDialogs } from '@toolpad/core/useDialogs'

import AddIcon from '@mui/icons-material/Add'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import DialogNewStreamConfig from '@/components/DialogNewStreamConfig'

type ButtonNewStreamConfigProps = Readonly<{
  onStreamConfigCreated: (name: string) => void
}>

const ButtonNewStreamConfig = (props: ButtonNewStreamConfigProps) => {
  const dialogs = useDialogs()
  const notify = useNotify()

  const [handleClick] = useApiOperation(async () => {
    const streamName = (await dialogs.open(DialogNewStreamConfig)) as
      | string
      | null
    if (streamName !== null) {
      props.onStreamConfigCreated(streamName)

      notify.success('Stream created')
    }
  }, [props.onStreamConfigCreated])

  return (
    <Button
      id="btn:streams.create"
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

export default ButtonNewStreamConfig
