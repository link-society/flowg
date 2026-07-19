import { useTranslation } from 'react-i18next'

import Button from '@mui/material/Button'

import AddIcon from '@mui/icons-material/Add'

import { useApiOperation } from '@/lib/hooks/api'
import { useDialogs } from '@/lib/hooks/dialogs'
import { useNotify } from '@/lib/hooks/notify'

import DialogNewTransformer from '@/components/DialogNewTransformer/component'

import { ButtonNewTransformerProps } from './types'

const ButtonNewTransformer = ({
  onTransformerCreated,
}: ButtonNewTransformerProps) => {
  const { t } = useTranslation()
  const dialogs = useDialogs()
  const notify = useNotify()

  const [handleClick] = useApiOperation(async () => {
    const transformerName = (await dialogs.open(DialogNewTransformer)) as
      string | null
    if (transformerName !== null) {
      onTransformerCreated(transformerName)

      notify.success(t('components.buttonNewTransformer.notifications.created'))
    }
  }, [onTransformerCreated])

  return (
    <Button
      id="btn:transformers.create"
      variant="contained"
      color="primary"
      size="small"
      onClick={() => handleClick()}
      startIcon={<AddIcon />}
    >
      {t('components.buttonNewTransformer.label')}
    </Button>
  )
}

export default ButtonNewTransformer
