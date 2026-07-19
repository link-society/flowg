import { useEffect } from 'react'
import { useTranslation } from 'react-i18next'

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
  const { t } = useTranslation()
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
        label={t('components.forwarderEditorAzureMonitor.endpointLabel')}
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
        label={t('components.forwarderEditorAzureMonitor.tokenLabel')}
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
        label={t('components.forwarderEditorAzureMonitor.expiresOnLabel')}
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
        label={t('components.forwarderEditorAzureMonitor.ruleIdLabel')}
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
        label={t('components.forwarderEditorAzureMonitor.streamNameLabel')}
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
          label={t('components.forwarderEditorAzureMonitor.allowInsecureLabel')}
        />
      </FormGroup>
    </ForwarderEditorAzureMonitorRoot>
  )
}

export default ForwarderEditorAwsCloudWatch
