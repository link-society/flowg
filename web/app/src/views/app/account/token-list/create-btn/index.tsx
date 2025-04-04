import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import { useDialogs } from '@toolpad/core/useDialogs'

import AddIcon from '@mui/icons-material/Add'

import * as tokenApi from '@/lib/api/operations/token'
import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import { ShowNewTokenModal } from './modal'

type CreateTokenButtonProps = Readonly<{
  onTokenCreated: (uuid: string) => void
}>

export const CreateTokenButton = ({
  onTokenCreated,
}: CreateTokenButtonProps) => {
  const dialogs = useDialogs()
  const notify = useNotify()

  const [handleClick, loading] = useApiOperation(async () => {
    const resp = await tokenApi.createToken()
    await dialogs.open(ShowNewTokenModal, resp)
    onTokenCreated(resp.token_uuid)
    notify.success('Token created')
  }, [onTokenCreated])

  return (
    <Button
      id="btn:account.tokens.create"
      variant="contained"
      color="secondary"
      size="small"
      disabled={loading}
      onClick={() => handleClick()}
    >
      {loading ? <CircularProgress color="inherit" size={24} /> : <AddIcon />}
    </Button>
  )
}
