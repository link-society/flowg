import InputAdornment from '@mui/material/InputAdornment'
import TextField from '@mui/material/TextField'

import { OpenTelemetryIcon } from '@/components/icons/otlp'

import { ForwarderTypeLabelMap } from '@/lib/models/forwarder'
import { OtlpForwarderModel } from '@/lib/models/forwarder/otlp'

import { KeyValueEditor } from '@/components/form/kv-editor'

type OtlpForwarderEditorProps = {
  config: OtlpForwarderModel
  onConfigChange: (config: OtlpForwarderModel) => void
}

export const OtlpForwarderEditor = ({
  config,
  onConfigChange,
}: OtlpForwarderEditorProps) => {
  const headerPairs = Object.entries(config.headers || {}).map(
    ([key, value]) => [key, value] as [string, string]
  )

  const onHeadersChange = (pairs: [string, string][]) => {
    const newHeaders = pairs.reduce(
      (acc, [key, value]) => ({
        ...acc,
        [key]: value,
      }),
      {} as Record<string, string>
    )

    // Only include headers if there are any
    const updatedConfig = {
      ...config,
      endpoint: config.endpoint,
    }
    if (Object.keys(newHeaders).length > 0) {
      updatedConfig.headers = newHeaders
    } else {
      delete updatedConfig.headers
    }

    onConfigChange(updatedConfig)
  }

  return (
    <div
      id="container:editor.forwarders.otlp"
      className="flex flex-col items-stretch gap-3"
    >
      <div className="mb-6 shadow">
        <TextField
          label="Forwarder Type"
          variant="outlined"
          className="w-full"
          type="text"
          value={ForwarderTypeLabelMap.otlp}
          disabled
          slotProps={{
            input: {
              startAdornment: (
                <InputAdornment position="start">
                  <OpenTelemetryIcon />
                </InputAdornment>
              ),
            },
          }}
        />
      </div>

      <TextField
        id="input:editor.forwarders.otlp.endpoint"
        label="Endpoint"
        value={config.endpoint}
        onChange={(e) =>
          onConfigChange({
            ...config,
            endpoint: e.target.value,
          })
        }
        type="text"
        variant="outlined"
        required
        className="w-full"
      />

      <KeyValueEditor
        id="input:editor.forwarders.otlp.headers"
        keyLabel="Header Name"
        valueLabel="Header Value"
        keyValues={headerPairs}
        onChange={onHeadersChange}
      />
    </div>
  )
}
