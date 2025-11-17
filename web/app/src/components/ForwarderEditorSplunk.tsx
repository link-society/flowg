import Divider from '@mui/material/Divider'
import InputAdornment from '@mui/material/InputAdornment'
import TextField from '@mui/material/TextField'

import { ForwarderConfigTypeLabelMap } from '@/lib/models/ForwarderConfigModel'
import ForwarderConfigSplunkModel from '@/lib/models/ForwarderConfigSplunkModel'

import ForwarderIconSplunk from '@/components/ForwarderIconSplunk'

type ForwarderEditorSplunkProps = {
  config: ForwarderConfigSplunkModel
  onConfigChange: (config: ForwarderConfigSplunkModel) => void
}

const ForwarderEditorSplunk = ({
  config,
  onConfigChange,
}: ForwarderEditorSplunkProps) => {
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
          value={ForwarderConfigTypeLabelMap.splunk}
          disabled
          slotProps={{
            input: {
              startAdornment: (
                <InputAdornment position="start">
                  <ForwarderIconSplunk />
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

export default ForwarderEditorSplunk
