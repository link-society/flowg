import { useDialogs } from '@toolpad/core/useDialogs'
import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

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
  const notify = useNotify()

  const [handleClick, loading] = useApiOperation(
    async () => {
      const { token, token_uuid } = await tokenApi.createToken()
      await dialogs.open(ShowNewTokenModal, token)
      onTokenCreated(token_uuid)
      notify.success('Token created')
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
