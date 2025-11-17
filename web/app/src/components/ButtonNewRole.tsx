import Button from '@mui/material/Button'
import { useDialogs } from '@toolpad/core/useDialogs'

import AddIcon from '@mui/icons-material/Add'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import RoleModel from '@/lib/models/RoleModel'

import DialogNewRole from '@/components/DialogNewRole'

type ButtonNewRoleProps = Readonly<{
  onRoleCreated: (role: RoleModel) => void
}>

const ButtonNewRole = ({ onRoleCreated }: ButtonNewRoleProps) => {
  const dialogs = useDialogs()
  const notify = useNotify()

  const [handleClick] = useApiOperation(async () => {
    const role = (await dialogs.open(DialogNewRole)) as RoleModel | null
    if (role !== null) {
      onRoleCreated(role)
      notify.success('Role created')
    }
  }, [onRoleCreated])

  return (
    <Button
      id="btn:admin.roles.create"
      variant="contained"
      color="secondary"
      size="small"
      onClick={() => handleClick()}
    >
      <AddIcon />
    </Button>
  )
}

export default ButtonNewRole
