import TextField from '@mui/material/TextField'

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
    <div className="flex flex-col gap-4">
      <TextField
        id="input:forwarder.otlp.endpoint"
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
        id="input:forwarder.otlp.headers"
        keyLabel="Header Name"
        valueLabel="Header Value"
        keyValues={headerPairs}
        onChange={onHeadersChange}
      />
    </div>
  )
}
