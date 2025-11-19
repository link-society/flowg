import { useEffect } from 'react'

import Divider from '@mui/material/Divider'
import TextField from '@mui/material/TextField'

import ForwarderConfigHttpModel from '@/lib/models/ForwarderConfigHttpModel'

import { useInput } from '@/lib/hooks/input'

import * as validators from '@/lib/validators'

import InputKeyValue from '@/components/InputKeyValue'

type ForwarderEditorHttpProps = {
  config: ForwarderConfigHttpModel
  onConfigChange: (config: ForwarderConfigHttpModel) => void
  onValidationChange: (valid: boolean) => void
}

const ForwarderEditorHttp = ({
  config,
  onConfigChange,
  onValidationChange,
}: ForwarderEditorHttpProps) => {
  const [url, setUrl] = useInput<string>(config.url, [
    validators.minLength(1),
    validators.formatUri,
  ])
  const [headers, setHeaders] = useInput<Record<string, string>>(
    config.headers ?? {},
    []
  )

  useEffect(() => {
    const valid = url.valid && headers.valid
    onValidationChange(valid)

    if (valid) {
      onConfigChange({
        type: 'http',
        url: url.value,
        headers: headers.value,
      })
    }
  }, [url, headers, onValidationChange, onConfigChange])

  return (
    <div
      id="container:editor.forwarders.http"
      className="flex flex-col items-stretch gap-3"
    >
      <TextField
        id="input:editor.forwarders.http.webhook_url"
        label="Webhook URL"
        variant="outlined"
        type="text"
        error={!url.valid}
        value={url.value}
        onChange={(e) => { setUrl(e.target.value) }}
      />

      <Divider />

      <InputKeyValue
        id="field:editor.forwarders.http.headers"
        keyLabel="HTTP Header"
        valueLabel="Value"
        keyValues={Object.entries(headers.value ?? {})}
        onChange={(pairs) => { setHeaders(Object.fromEntries(pairs)) }}
      />
    </div>
  )
}

export default ForwarderEditorHttp
