import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import { useDialogs } from '@toolpad/core/useDialogs'

import AddIcon from '@mui/icons-material/Add'

import * as tokenApi from '@/lib/api/operations/token'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import DialogNewToken from '@/components/DialogNewToken'

type ButtonNewTokenProps = Readonly<{
  onTokenCreated: (uuid: string) => void
}>

const ButtonNewToken = ({ onTokenCreated }: ButtonNewTokenProps) => {
  const dialogs = useDialogs()
  const notify = useNotify()

  const [handleClick, loading] = useApiOperation(async () => {
    const resp = await tokenApi.createToken()
    await dialogs.open(DialogNewToken, resp)
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

export default ButtonNewToken
