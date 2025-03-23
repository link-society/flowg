import { useDialogs } from '@toolpad/core/useDialogs'
import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import AddIcon from '@mui/icons-material/Add'

import Button from '@mui/material/Button'

import { UserModel } from '@/lib/models/auth'

import { UserFormModal } from './modal'

type CreateUserButtonProps = Readonly<{
  roles: string[]
  onUserCreated: (user: UserModel) => void
}>

export const CreateUserButton = ({ roles, onUserCreated }: CreateUserButtonProps) => {
  const dialogs = useDialogs()
  const notify = useNotify()

  const [handleClick] = useApiOperation(
    async () => {
      const user = await dialogs.open(UserFormModal, roles) as UserModel | null
      if (user !== null) {
        onUserCreated(user)
        notify.success('User created')
      }
    },
    [onUserCreated],
  )

  return (
    <Button
      id="btn:admin.users.create"
      variant="contained"
      color="secondary"
      size="small"
      onClick={() => handleClick()}
    >
      <AddIcon />
    </Button>
  )
}
