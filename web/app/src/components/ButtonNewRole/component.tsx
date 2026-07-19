import { useTranslation } from 'react-i18next'

import Button from '@mui/material/Button'

import AddIcon from '@mui/icons-material/Add'

import { useApiOperation } from '@/lib/hooks/api'
import { useDialogs } from '@/lib/hooks/dialogs'
import { useNotify } from '@/lib/hooks/notify'

import RoleModel from '@/lib/models/RoleModel'

import DialogNewRole from '@/components/DialogNewRole/component'

import { ButtonNewRoleProps } from './types'

const ButtonNewRole = ({ onRoleCreated }: ButtonNewRoleProps) => {
  const { t } = useTranslation()
  const dialogs = useDialogs()
  const notify = useNotify()

  const [handleClick] = useApiOperation(async () => {
    const role = (await dialogs.open(DialogNewRole)) as RoleModel | null
    if (role !== null) {
      onRoleCreated(role)
      notify.success(t('components.buttonNewRole.notifications.created'))
    }
  }, [onRoleCreated])

  return (
    <Button
      id="btn:admin.roles.create"
      variant="contained"
      size="small"
      color="secondary"
      onClick={() => handleClick()}
    >
      <AddIcon />
    </Button>
  )
}

export default ButtonNewRole
