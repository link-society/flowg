import Divider from '@mui/material/Divider'
import InputAdornment from '@mui/material/InputAdornment'
import TextField from '@mui/material/TextField'

import ForwarderConfigDatadogModel from '@/lib/models/ForwarderConfigDatadogModel'
import { ForwarderConfigTypeLabelMap } from '@/lib/models/ForwarderConfigModel'

import ForwarderIconDatadog from '@/components/ForwarderIconDatadog'

type ForwarderEditorDatadogProps = {
  config: ForwarderConfigDatadogModel
  onConfigChange: (config: ForwarderConfigDatadogModel) => void
}

const ForwarderEditorDatadog = ({
  config,
  onConfigChange,
}: ForwarderEditorDatadogProps) => {
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
          value={ForwarderConfigTypeLabelMap.datadog}
          disabled
          slotProps={{
            input: {
              startAdornment: (
                <InputAdornment position="start">
                  <ForwarderIconDatadog />
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

export default ForwarderEditorDatadog
