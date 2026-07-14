import { useEffect } from 'react'

import Checkbox from '@mui/material/Checkbox'
import Divider from '@mui/material/Divider'
import FormControlLabel from '@mui/material/FormControlLabel'
import FormGroup from '@mui/material/FormGroup'
import TextField from '@mui/material/TextField'

import { useInput } from '@/lib/hooks/input'

import * as validators from '@/lib/validators'

import { ForwarderEditorAzureMonitorRoot } from './styles'
import { ForwarderEditorAzureMonitorProps } from './types'

const ForwarderEditorAwsCloudWatch = ({
  config,
  onConfigChange,
  onValidationChange,
}: ForwarderEditorAzureMonitorProps) => {
  const [endpoint, setEndpoint] = useInput(config.endpoint, [
    validators.minLength(1),
    validators.formatUri,
  ])
  const [token, setToken] = useInput(config.token, [validators.minLength(1)])
  const [expires_on, setExpiresOn] = useInput(config.expires_on, [
    validators.minLength(1),
  ])
  const [rule_id, setRuleID] = useInput(config.rule_id, [
    validators.minLength(1),
  ])
  const [stream_name, setStreamName] = useInput(config.stream_name, [
    validators.minLength(1),
  ])
  const [allow_insecure, setAllowInsecure] = useInput(config.allow_insecure)

  useEffect(() => {
    const valid = true
    onValidationChange(valid)

    if (valid) {
      onConfigChange({
        type: 'azuremonitor',
        endpoint: endpoint.value,
        token: token.value,
        expires_on: expires_on.value,
        rule_id: rule_id.value,
        stream_name: stream_name.value,
        allow_insecure: allow_insecure.value,
      })
    }
  }, [endpoint, token, expires_on, rule_id, stream_name, allow_insecure])

  return (
    <ForwarderEditorAzureMonitorRoot id="container:editor.forwarders.azuremonitor">
      <TextField
        id="input:editor.forwarders.azuremonitor.endpoint"
        label="Endpoint"
        variant="outlined"
        type="text"
        error={!endpoint.valid}
        value={endpoint.value}
        onChange={(e) => {
          setEndpoint(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.azuremonitor.token"
        label="Token"
        variant="outlined"
        type="text"
        error={!token.valid}
        value={token.value}
        onChange={(e) => {
          setToken(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.azuremonitor.expires_on"
        label="Expires on"
        variant="outlined"
        type="text"
        error={!expires_on.valid}
        value={expires_on.value}
        onChange={(e) => {
          setExpiresOn(e.target.value)
        }}
      />

      <Divider />

      <TextField
        id="input:editor.forwarders.azuremonitor.rule_id"
        label="Rule ID"
        variant="outlined"
        type="text"
        error={!rule_id.valid}
        value={rule_id.value}
        onChange={(e) => {
          setRuleID(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.azuremonitor.stream_name"
        label="Stream name"
        variant="outlined"
        type="text"
        error={!stream_name.valid}
        value={stream_name.value}
        onChange={(e) => {
          setStreamName(e.target.value)
        }}
      />

      <FormGroup>
        <FormControlLabel
          control={
            <Checkbox
              id="input:editor.forwarders.azuremonitor.allow_insecure"
              checked={allow_insecure.value}
              onChange={(e) => {
                setAllowInsecure(e.target.checked)
              }}
            />
          }
          label="Allow insecure connections"
        />
      </FormGroup>
    </ForwarderEditorAzureMonitorRoot>
  )
}

export default ForwarderEditorAwsCloudWatch
