import { useEffect } from 'react'

import TextField from '@mui/material/TextField'

import { useInput } from '@/lib/hooks/input'

import ForwarderConfigOtlpModel from '@/lib/models/ForwarderConfigOtlpModel'

import * as validators from '@/lib/validators'

import InputKeyValue from '@/components/InputKeyValue'

type ForwarderEditorOtlpProps = {
  config: ForwarderConfigOtlpModel
  onConfigChange: (config: ForwarderConfigOtlpModel) => void
  onValidationChange: (valid: boolean) => void
}

const ForwarderEditorOtlp = ({
  config,
  onConfigChange,
  onValidationChange,
}: ForwarderEditorOtlpProps) => {
  const [endpoint, setEndpoint] = useInput<string>(config.endpoint, [
    validators.minLength(1),
    validators.formatUri,
  ])
  const [headers, setHeaders] = useInput<Record<string, string>>(
    config.headers ?? {},
    []
  )

  useEffect(() => {
    const valid = endpoint.valid && headers.valid
    onValidationChange(valid)

    if (valid) {
      onConfigChange({
        type: 'otlp',
        endpoint: endpoint.value,
        headers: headers.value,
      })
    }
  }, [endpoint, headers])

  return (
    <div
      id="container:editor.forwarders.otlp"
      className="flex flex-col items-stretch gap-3"
    >
      <TextField
        id="input:editor.forwarders.otlp.endpoint"
        label="Endpoint"
        error={!endpoint.valid}
        value={endpoint.value}
        onChange={(e) => setEndpoint(e.target.value)}
        type="text"
        variant="outlined"
        required
        className="w-full"
      />

      <InputKeyValue
        id="input:editor.forwarders.otlp.headers"
        keyLabel="Header Name"
        valueLabel="Header Value"
        keyValues={Object.entries(headers.value ?? {})}
        onChange={(pairs) => {
          setHeaders(Object.fromEntries(pairs))
        }}
      />
    </div>
  )
}

export default ForwarderEditorOtlp
