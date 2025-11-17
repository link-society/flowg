import Button from '@mui/material/Button'
import { useDialogs } from '@toolpad/core/useDialogs'

import AddIcon from '@mui/icons-material/Add'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import UserModel from '@/lib/models/UserModel'

import DialogNewUser from '@/components/DialogNewUser'

type ButtonNewUserProps = Readonly<{
  roles: string[]
  onUserCreated: (user: UserModel) => void
}>

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
