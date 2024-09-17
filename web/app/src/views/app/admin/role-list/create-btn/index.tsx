import { useDialogs } from '@toolpad/core/useDialogs'
import { useNavigate } from 'react-router-dom'
import { useNotifications } from '@toolpad/core/useNotifications'
import { useConfig } from '@/lib/context/config'

import AddIcon from '@mui/icons-material/Add'

import Button from '@mui/material/Button'

import { RoleModel } from '@/lib/models'
import { UnauthenticatedError } from '@/lib/api/errors'

import { RoleFormModal } from './modal'

type CreateRoleButtonProps = {
  onRoleCreated: (role: RoleModel) => void
}

export const CreateRoleButton = ({ onRoleCreated }: CreateRoleButtonProps) => {
  const dialogs = useDialogs()
  const navigate = useNavigate()
  const notifications = useNotifications()
  const config = useConfig()

  const handleClick = async () => {
    try {
      const role = await dialogs.open(RoleFormModal) as RoleModel | null
      if (role !== null) {
        onRoleCreated(role)
        notifications.show('Role created', {
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
