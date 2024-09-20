import { useDialogs } from '@toolpad/core/useDialogs'
import { useNotifications } from '@toolpad/core/useNotifications'
import { useConfig } from '@/lib/context/config'
import { useApiOperation } from '@/lib/hooks/api'

import AddIcon from '@mui/icons-material/Add'

import Button from '@mui/material/Button'

import { RoleModel } from '@/lib/models'

import { RoleFormModal } from './modal'

type CreateRoleButtonProps = {
  onRoleCreated: (role: RoleModel) => void
}

export const CreateRoleButton = ({ onRoleCreated }: CreateRoleButtonProps) => {
  const dialogs = useDialogs()
  const notifications = useNotifications()
  const config = useConfig()

  const [handleClick] = useApiOperation(
    async () => {
      const role = await dialogs.open(RoleFormModal) as RoleModel | null
      if (role !== null) {
        onRoleCreated(role)
        notifications.show('Role created', {
          severity: 'success',
          autoHideDuration: config.notifications?.autoHideDuration,
        })
      }
    },
    [onRoleCreated],
  )

  return (
    <Button
      variant="contained"
      color="secondary"
      size="small"
      onClick={() => handleClick()}
    >
      <AddIcon />
    </Button>
  )
}
