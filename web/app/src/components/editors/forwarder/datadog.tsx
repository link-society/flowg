import Divider from '@mui/material/Divider'
import InputAdornment from '@mui/material/InputAdornment'
import TextField from '@mui/material/TextField'

import { DatadogIcon } from '@/components/icons/datadog'

import { ForwarderTypeLabelMap } from '@/lib/models/forwarder'
import { DatadogForwarderModel } from '@/lib/models/forwarder/datadog'

type DatadogForwarderEditorProps = {
  config: DatadogForwarderModel
  onConfigChange: (config: DatadogForwarderModel) => void
}

export const DatadogForwarderEditor = ({
  config,
  onConfigChange,
}: DatadogForwarderEditorProps) => {
  return (
    <div
      id="container:editor.forwarders.datadog"
      className="flex flex-col items-stretch gap-3"
    >
      <div className="mb-6 shadow">
        <TextField
          label="Forwarder Type"
          variant="outlined"
          className="w-full"
          type="text"
          value={ForwarderTypeLabelMap.datadog}
          disabled
          slotProps={{
            input: {
              startAdornment: (
                <InputAdornment position="start">
                  <DatadogIcon />
                </InputAdornment>
              ),
            },
          }}
        />
      </div>

      <TextField
        id="input:editor.forwarders.datadog.url"
        label="URL"
        variant="outlined"
        type="text"
        value={config.url}
        onChange={(e) => {
          onConfigChange({
            ...config,
            url: e.target.value,
          })
        }}
      />

      <Divider />

      <TextField
        id="input:editor.forwarders.datadog.apiKey"
        label="ApiKey"
        variant="outlined"
        type="password"
        value={config.apiKey}
        onChange={(e) => {
          onConfigChange({
            ...config,
            apiKey: e.target.value,
          })
        }}
      />
    </div>
  )
}
