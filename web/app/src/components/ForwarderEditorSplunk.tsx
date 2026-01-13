import { useEffect } from 'react'

import Divider from '@mui/material/Divider'
import TextField from '@mui/material/TextField'

import { useInput } from '@/lib/hooks/input'

import ForwarderConfigSplunkModel from '@/lib/models/ForwarderConfigSplunkModel'

import * as validators from '@/lib/validators'

type ForwarderEditorSplunkProps = {
  config: ForwarderConfigSplunkModel
  onConfigChange: (config: ForwarderConfigSplunkModel) => void
  onValidationChange: (valid: boolean) => void
}

const ForwarderEditorSplunk = ({
  config,
  onConfigChange,
  onValidationChange,
}: ForwarderEditorSplunkProps) => {
  const [endpoint, setEndpoint] = useInput<string>(config.endpoint, [
    validators.minLength(1),
    validators.formatUri,
  ])
  const [token, setToken] = useInput<string>(config.token, [
    validators.minLength(1),
  ])

  useEffect(() => {
    const valid = endpoint.valid && token.valid
    onValidationChange(valid)

    if (valid) {
      onConfigChange({
        type: 'splunk',
        endpoint: endpoint.value,
        token: token.value,
      })
    }
  }, [endpoint, token])

  return (
    <div
      id="container:editor.forwarders.splunk"
      className="flex flex-col items-stretch gap-3"
    >
      <TextField
        id="input:editor.forwarders.splunk.endpoint"
        label="HTTP Event Collector Endpoint"
        variant="outlined"
        type="text"
        error={!endpoint.valid}
        value={endpoint.value}
        onChange={(e) => {
          setEndpoint(e.target.value)
        }}
      />

      <Divider />

      <TextField
        id="input:editor.forwarders.splunk.token"
        label="Token"
        variant="outlined"
        type="password"
        error={!token.valid}
        value={token.value}
        onChange={(e) => {
          setToken(e.target.value)
        }}
      />
    </div>
  )
}

export default ForwarderEditorSplunk
