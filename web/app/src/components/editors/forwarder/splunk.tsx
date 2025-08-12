import Divider from '@mui/material/Divider'
import InputAdornment from '@mui/material/InputAdornment'
import TextField from '@mui/material/TextField'

import { ForwarderTypeLabelMap } from '@/lib/models/forwarder'
import { SplunkForwarderModel } from '@/lib/models/forwarder/splunk'

import { SplunkIcon } from '@/components/icons/splunk'

type SplunkForwarderEditorProps = {
  config: SplunkForwarderModel
  onConfigChange: (config: SplunkForwarderModel) => void
}

export const SplunkForwarderEditor = ({
  config,
  onConfigChange,
}: SplunkForwarderEditorProps) => {
  return (
    <div
      id="container:editor.forwarders.splunk"
      className="flex flex-col items-stretch gap-3"
    >
      <div className="mb-6 shadow">
        <TextField
          label="Forwarder Type"
          variant="outlined"
          className="w-full"
          type="text"
          value={ForwarderTypeLabelMap.splunk}
          disabled
          slotProps={{
            input: {
              startAdornment: (
                <InputAdornment position="start">
                  <SplunkIcon />
                </InputAdornment>
              ),
            },
          }}
        />
      </div>

      <TextField
        id="input:editor.forwarders.splunk.endpoint"
        label="HTTP Event Collector Endpoint"
        variant="outlined"
        type="text"
        value={config.endpoint}
        onChange={(e) => {
          onConfigChange({
            ...config,
            endpoint: e.target.value,
          })
        }}
      />

      <Divider />

      <TextField
        id="input:editor.forwarders.splunk.token"
        label="Token"
        variant="outlined"
        type="password"
        value={config.token}
        onChange={(e) => {
          onConfigChange({
            ...config,
            token: e.target.value,
          })
        }}
      />
    </div>
  )
}
