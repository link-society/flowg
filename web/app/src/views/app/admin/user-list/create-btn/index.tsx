import { useDialogs } from '@toolpad/core/useDialogs'
import { useNavigate } from 'react-router-dom'
import { useNotifications } from '@toolpad/core/useNotifications'
import { useConfig } from '@/lib/context/config'

import AddIcon from '@mui/icons-material/Add'

import Button from '@mui/material/Button'

import { UserModel } from '@/lib/models'
import { UnauthenticatedError } from '@/lib/api/errors'

import { UserFormModal } from './modal'

type CreateUserButtonProps = {
  roles: string[]
  onUserCreated: (user: UserModel) => void
}

export const CreateUserButton = ({ roles, onUserCreated }: CreateUserButtonProps) => {
  const dialogs = useDialogs()
  const navigate = useNavigate()
  const notifications = useNotifications()
  const config = useConfig()

  const handleClick = async () => {
    try {
      const user = await dialogs.open(UserFormModal, roles) as UserModel | null
      if (user !== null) {
        onUserCreated(user)
        notifications.show('User created', {
          severity: 'success',
          autoHideDuration: config.notifications?.autoHideDuration,
        })
      }
    }
    catch (error) {
      if (error instanceof UnauthenticatedError) {
        notifications.show('Session expired', {
          severity: 'error',
          autoHideDuration: config.notifications?.autoHideDuration,
        })
        navigate('/web/login')
      }
      else {
        notifications.show('Unknown error', {
          severity: 'error',
          autoHideDuration: config.notifications?.autoHideDuration,
        })
      }

      console.error(error)
    }
  }

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
