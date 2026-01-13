import { useEffect } from 'react'

import Divider from '@mui/material/Divider'
import TextField from '@mui/material/TextField'

import { useInput } from '@/lib/hooks/input'

import ForwarderConfigDatadogModel from '@/lib/models/ForwarderConfigDatadogModel'

import * as validators from '@/lib/validators'

type ForwarderEditorDatadogProps = {
  config: ForwarderConfigDatadogModel
  onConfigChange: (config: ForwarderConfigDatadogModel) => void
  onValidationChange: (valid: boolean) => void
}

const ForwarderEditorDatadog = ({
  config,
  onConfigChange,
  onValidationChange,
}: ForwarderEditorDatadogProps) => {
  const [url, setUrl] = useInput(config.url, [
    validators.minLength(1),
    validators.formatUri,
  ])
  const [apiKey, setApiKey] = useInput(config.apiKey, [validators.minLength(1)])

  useEffect(() => {
    const valid = url.valid && apiKey.valid
    onValidationChange(valid)

    if (valid) {
      onConfigChange({
        type: 'datadog',
        url: url.value,
        apiKey: apiKey.value,
      })
    }
  }, [url, apiKey])

  return (
    <div
      id="container:editor.forwarders.datadog"
      className="flex flex-col items-stretch gap-3"
    >
      <TextField
        id="input:editor.forwarders.datadog.url"
        label="URL"
        variant="outlined"
        type="text"
        error={!url.valid}
        value={url.value}
        onChange={(e) => {
          setUrl(e.target.value)
        }}
      />

      <Divider />

      <TextField
        id="input:editor.forwarders.datadog.apiKey"
        label="ApiKey"
        variant="outlined"
        type="password"
        error={!apiKey.valid}
        value={apiKey.value}
        onChange={(e) => {
          setApiKey(e.target.value)
        }}
      />
    </div>
  )
}

export default ForwarderEditorDatadog
