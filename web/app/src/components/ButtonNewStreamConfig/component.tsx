import { useTranslation } from 'react-i18next'

import Button from '@mui/material/Button'

import AddIcon from '@mui/icons-material/Add'

import { useApiOperation } from '@/lib/hooks/api'
import { useDialogs } from '@/lib/hooks/dialogs'
import { useNotify } from '@/lib/hooks/notify'

import DialogNewStreamConfig from '@/components/DialogNewStreamConfig/component'

import { ButtonNewStreamConfigProps } from './types'

const ButtonNewStreamConfig = (props: ButtonNewStreamConfigProps) => {
  const { t } = useTranslation()
  const dialogs = useDialogs()
  const notify = useNotify()

  const [handleClick] = useApiOperation(async () => {
    const streamName = (await dialogs.open(DialogNewStreamConfig)) as
      string | null
    if (streamName !== null) {
      props.onStreamConfigCreated(streamName)

      notify.success(
        t('components.buttonNewStreamConfig.notifications.created')
      )
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
      {t('components.buttonNewStreamConfig.label')}
    </Button>
  )
}

export default ButtonNewStreamConfig
