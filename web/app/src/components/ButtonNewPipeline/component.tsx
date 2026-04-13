import Button from '@mui/material/Button'

import AddIcon from '@mui/icons-material/Add'

import { useApiOperation } from '@/lib/hooks/api'
import { useDialogs } from '@/lib/hooks/dialogs'
import { useNotify } from '@/lib/hooks/notify'

import DialogNewPipeline from '@/components/DialogNewPipeline/component'

import { ButtonNewPipelineProps } from './types'

const ButtonNewPipeline = (props: ButtonNewPipelineProps) => {
  const dialogs = useDialogs()
  const notify = useNotify()

  const [handleClick] = useApiOperation(async () => {
    const pipelineName = (await dialogs.open(DialogNewPipeline)) as
      | string
      | null
    if (pipelineName !== null) {
      props.onPipelineCreated(pipelineName)

      notify.success('Pipeline created')
    }
  }, [props.onPipelineCreated])

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

export default ButtonNewPipeline
