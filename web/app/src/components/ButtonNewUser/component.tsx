import { useTranslation } from 'react-i18next'

import Button from '@mui/material/Button'

import AddIcon from '@mui/icons-material/Add'

import { useApiOperation } from '@/lib/hooks/api'
import { useDialogs } from '@/lib/hooks/dialogs'
import { useNotify } from '@/lib/hooks/notify'

import UserModel from '@/lib/models/UserModel'

import DialogNewUser from '@/components/DialogNewUser/component'

import { ButtonNewUserProps } from './types'

const ButtonNewUser = ({
  roles,
  defaultRoles,
  onUserCreated,
}: ButtonNewUserProps) => {
  const { t } = useTranslation()
  const dialogs = useDialogs()
  const notify = useNotify()

  const [handleClick] = useApiOperation(async () => {
    const user = (await dialogs.open(DialogNewUser, {
      roles,
      defaultRoles,
    })) as UserModel | null
    if (user !== null) {
      onUserCreated(user)
      notify.success(t('components.buttonNewUser.notifications.created'))
    }
  }, [onUserCreated])

  return (
    <Button
      id="btn:admin.users.create"
      variant="contained"
      size="small"
      color="secondary"
      onClick={() => handleClick()}
    >
      <AddIcon />
    </Button>
  )
}

export default ButtonNewUser
