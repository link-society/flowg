import { useDialogs } from '@toolpad/core/useDialogs'
import { useNotifications } from '@toolpad/core/useNotifications'
import { useConfig } from '@/lib/context/config'
import { useApiOperation } from '@/lib/hooks/api'

import AddIcon from '@mui/icons-material/Add'

import Button from '@mui/material/Button'

import { UserModel } from '@/lib/models'

import { UserFormModal } from './modal'

type CreateUserButtonProps = {
  roles: string[]
  onUserCreated: (user: UserModel) => void
}

export const CreateUserButton = ({ roles, onUserCreated }: CreateUserButtonProps) => {
  const dialogs = useDialogs()
  const notifications = useNotifications()
  const config = useConfig()

  const [handleClick] = useApiOperation(
    async () => {
      const user = await dialogs.open(UserFormModal, roles) as UserModel | null
      if (user !== null) {
        onUserCreated(user)
        notifications.show('User created', {
          severity: 'success',
          autoHideDuration: config.notifications?.autoHideDuration,
        })
      }
    },
    [onUserCreated],
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
