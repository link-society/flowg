import { useDialogs } from '@toolpad/core/useDialogs'
import { useNotifications } from '@toolpad/core/useNotifications'
import { useConfig } from '@/lib/context/config'
import { useApiOperation } from '@/lib/hooks/api'

import AddIcon from '@mui/icons-material/Add'

import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'

import * as tokenApi from '@/lib/api/operations/token'

import { ShowNewTokenModal } from './modal'

type CreateTokenButtonProps = {
  onTokenCreated: (uuid: string) => void
}

export const CreateTokenButton = ({ onTokenCreated }: CreateTokenButtonProps) => {
  const dialogs = useDialogs()
  const notifications = useNotifications()
  const config = useConfig()

  const [handleClick, loading] = useApiOperation(
    async () => {
      const { token, token_uuid } = await tokenApi.createToken()
      await dialogs.open(ShowNewTokenModal, token)
      onTokenCreated(token_uuid)
      notifications.show('Token created', {
        severity: 'success',
        autoHideDuration: config.notifications?.autoHideDuration,
      })
    },
    [onTokenCreated],
  )

  return (
    <Button
      variant="contained"
      color="secondary"
      size="small"
      disabled={loading}
      onClick={() => handleClick()}
    >
      {loading
        ? <CircularProgress color="inherit" size={24} />
        : <AddIcon />
      }
    </Button>
  )
}
