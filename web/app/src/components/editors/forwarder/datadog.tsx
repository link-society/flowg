import { DatadogForwarderModel } from '@/lib/models/forwarder/datadog'
import Divider from '@mui/material/Divider'
import { ForwarderTypeLabelMap } from '@/lib/models/forwarder'
import HttpIcon from '@mui/icons-material/Http'
import InputAdornment from '@mui/material/InputAdornment'
import TextField from '@mui/material/TextField'

type DatadogForwarderEditorProps = {
  config: DatadogForwarderModel
  onConfigChange: (config: DatadogForwarderModel) => void
}

export const DatadogForwarderEditor = ({ config, onConfigChange }: DatadogForwarderEditorProps) => {
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
                  <HttpIcon />
                </InputAdornment>
              ),
            }
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
        type="text"
        value={config.apiKey}
        onChange={(e) => {
          onConfigChange({
            ...config,
            apiKey: e.target.value,
          })
        }}
      />

    <Divider /> 

    <TextField
        id="input:editor.forwarders.datadog.source"
        label="Source"
        variant="outlined"
        type="text"
        value={config.apiKey}
        onChange={(e) => {
          onConfigChange({
            ...config,
            source: e.target.value,
          })
        }}
      />

    <Divider />

    <TextField
        id="input:editor.forwarders.datadog.service"
        label="Service"
        variant="outlined"
        type="text"
        value={config.apiKey}
        onChange={(e) => {
          onConfigChange({
            ...config,
            service: e.target.value,
          })
        }}
      />
    </div>
    
  )
}
