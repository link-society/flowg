import { useDialogs } from '@toolpad/core/useDialogs'
import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import AddIcon from '@mui/icons-material/Add'

import Button from '@mui/material/Button'

import { RoleModel } from '@/lib/models'

import { RoleFormModal } from './modal'

type CreateRoleButtonProps = {
  onRoleCreated: (role: RoleModel) => void
}

export const CreateRoleButton = ({ onRoleCreated }: CreateRoleButtonProps) => {
  const dialogs = useDialogs()
  const notify = useNotify()

  const [handleClick] = useApiOperation(
    async () => {
      const role = await dialogs.open(RoleFormModal) as RoleModel | null
      if (role !== null) {
        onRoleCreated(role)
        notify.success('Role created')
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
