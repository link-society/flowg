import InputAdornment from '@mui/material/InputAdornment'
import TextField from '@mui/material/TextField'

import { ForwarderConfigTypeLabelMap } from '@/lib/models/ForwarderConfigModel'
import ForwarderConfigOtlpModel from '@/lib/models/ForwarderConfigOtlpModel'

import ForwarderIconOtlp from '@/components/ForwarderIconOtlp'
import InputKeyValue from '@/components/InputKeyValue'

type ForwarderEditorOtlpProps = {
  config: ForwarderConfigOtlpModel
  onConfigChange: (config: ForwarderConfigOtlpModel) => void
}

const ForwarderEditorOtlp = ({
  config,
  onConfigChange,
}: ForwarderEditorOtlpProps) => {
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
          value={ForwarderConfigTypeLabelMap.otlp}
          disabled
          slotProps={{
            input: {
              startAdornment: (
                <InputAdornment position="start">
                  <ForwarderIconOtlp />
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

      <InputKeyValue
        id="input:editor.forwarders.otlp.headers"
        keyLabel="Header Name"
        valueLabel="Header Value"
        keyValues={headerPairs}
        onChange={onHeadersChange}
      />
    </div>
  )
}

export default ForwarderEditorOtlp
