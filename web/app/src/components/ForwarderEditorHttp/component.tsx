import { useEffect } from 'react'

import Divider from '@mui/material/Divider'
import TextField from '@mui/material/TextField'

import { useInput } from '@/lib/hooks/input'

import { DynamicField } from '@/lib/models/DynamicField.ts'

import * as validators from '@/lib/validators'

import DynamicFieldControl from '@/components/DynamicFieldControl/component'
import InputKeyValue from '@/components/InputKeyValue/component'

import { ForwarderEditorHttpRoot } from './styles'
import { ForwarderEditorHttpProps } from './types'

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

  const [body, setBody] = useInput<DynamicField<string>>(config.body, [
    validators.dynamicField([]),
  ])

  useEffect(() => {
    const valid = url.valid && headers.valid
    onValidationChange(valid)

    if (valid) {
      onConfigChange({
        type: 'http',
        url: url.value,
        headers: headers.value,
        body: body.value,
      })
    }
  }, [url, headers])

  return (
    <ForwarderEditorHttpRoot id="container:editor.forwarders.http">
      <TextField
        id="input:editor.forwarders.http.webhook_url"
        label="Webhook URL"
        variant="outlined"
        type="text"
        error={!url.valid}
        value={url.value}
        onChange={(e) => {
          setUrl(e.target.value)
        }}
      />

      <Divider />

      <InputKeyValue
        id="field:editor.forwarders.http.headers"
        keyLabel="HTTP Header"
        valueLabel="Value"
        keyValues={Object.entries(headers.value ?? {})}
        onChange={(pairs) => {
          setHeaders(Object.fromEntries(pairs))
        }}
      />

      <DynamicFieldControl
        id="input:editor.forwarders.http.body"
        label="Body"
        multiline
        variant="outlined"
        error={!body.valid}
        value={body.value}
        onChange={setBody}
      />
    </ForwarderEditorHttpRoot>
  )
}

export default ForwarderEditorHttp
