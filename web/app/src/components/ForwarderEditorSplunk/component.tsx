import { useEffect } from 'react'
import { useTranslation } from 'react-i18next'

import Divider from '@mui/material/Divider'
import TextField from '@mui/material/TextField'

import { useInput } from '@/lib/hooks/input'

import { DynamicField } from '@/lib/models/DynamicField.ts'

import * as validators from '@/lib/validators'

import DynamicFieldControl from '@/components/DynamicFieldControl/component'

import { ForwarderEditorSplunkRoot } from './styles'
import { ForwarderEditorSplunkProps } from './types'

const ForwarderEditorSplunk = ({
  config,
  onConfigChange,
  onValidationChange,
}: ForwarderEditorSplunkProps) => {
  const { t } = useTranslation()
  const [endpoint, setEndpoint] = useInput<string>(config.endpoint, [
    validators.minLength(1),
    validators.formatUri,
  ])
  const [token, setToken] = useInput<string>(config.token, [
    validators.minLength(1),
  ])
  const [source, setSource] = useInput<DynamicField<string>>(config.source, [
    validators.dynamicField([]),
  ])
  const [host, setHost] = useInput<DynamicField<string>>(config.host, [
    validators.dynamicField([]),
  ])

  useEffect(() => {
    const valid = endpoint.valid && token.valid
    onValidationChange(valid)

    if (valid) {
      onConfigChange({
        type: 'splunk',
        endpoint: endpoint.value,
        token: token.value,
        source: source.value,
        host: host.value,
      })
    }
  }, [endpoint, token])

  return (
    <ForwarderEditorSplunkRoot id="container:editor.forwarders.splunk">
      <TextField
        id="input:editor.forwarders.splunk.endpoint"
        label={t('components.forwarderEditorSplunk.endpointLabel')}
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
        label={t('components.forwarderEditorSplunk.tokenLabel')}
        variant="outlined"
        type="password"
        error={!token.valid}
        value={token.value}
        onChange={(e) => {
          setToken(e.target.value)
        }}
      />

      <DynamicFieldControl
        id="input:editor.forwarders.splunk.source"
        label={t('components.forwarderEditorSplunk.sourceLabel')}
        variant="outlined"
        type="text"
        error={!source.valid}
        value={source.value}
        onChange={setSource}
      />

      <DynamicFieldControl
        id="input:editor.forwarders.splunk.host"
        label={t('components.forwarderEditorSplunk.hostLabel')}
        variant="outlined"
        type="text"
        error={!host.valid}
        value={host.value}
        onChange={setHost}
      />
    </ForwarderEditorSplunkRoot>
  )
}

export default ForwarderEditorSplunk
