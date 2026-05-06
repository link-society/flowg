import Button from '@mui/material/Button'

import AddIcon from '@mui/icons-material/Add'

import { useApiOperation } from '@/lib/hooks/api'
import { useDialogs } from '@/lib/hooks/dialogs'
import { useNotify } from '@/lib/hooks/notify'

import UserModel from '@/lib/models/UserModel'

import DialogNewUser from '@/components/DialogNewUser/component'

import { ButtonNewUserProps } from './types'

const ButtonNewUser = ({ roles, onUserCreated }: ButtonNewUserProps) => {
  const dialogs = useDialogs()
  const notify = useNotify()

  const [handleClick] = useApiOperation(async () => {
    const user = (await dialogs.open(DialogNewUser, roles)) as UserModel | null
    if (user !== null) {
      onUserCreated(user)
      notify.success('User created')
    }
  }, [onUserCreated])

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

export default ButtonNewUser
