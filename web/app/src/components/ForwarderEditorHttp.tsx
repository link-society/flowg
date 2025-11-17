import Divider from '@mui/material/Divider'
import InputAdornment from '@mui/material/InputAdornment'
import TextField from '@mui/material/TextField'

import HttpIcon from '@mui/icons-material/Http'

import ForwarderConfigHttpModel from '@/lib/models/ForwarderConfigHttpModel'
import { ForwarderConfigTypeLabelMap } from '@/lib/models/ForwarderConfigModel'

import InputKeyValue from '@/components/InputKeyValue'

type ForwarderEditorHttpProps = {
  config: ForwarderConfigHttpModel
  onConfigChange: (config: ForwarderConfigHttpModel) => void
}

const ForwarderEditorHttp = ({
  config,
  onConfigChange,
}: ForwarderEditorHttpProps) => {
  return (
    <div
      id="container:editor.forwarders.http"
      className="flex flex-col items-stretch gap-3"
    >
      <div className="mb-6 shadow">
        <TextField
          label="Forwarder Type"
          variant="outlined"
          className="w-full"
          type="text"
          value={ForwarderConfigTypeLabelMap.http}
          disabled
          slotProps={{
            input: {
              startAdornment: (
                <InputAdornment position="start">
                  <HttpIcon />
                </InputAdornment>
              ),
            },
          }}
        />
      </div>

      <TextField
        id="input:editor.forwarders.http.webhook_url"
        label="Webhook URL"
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

      <InputKeyValue
        id="field:editor.forwarders.http.headers"
        keyLabel="HTTP Header"
        valueLabel="Value"
        keyValues={Object.entries(config.headers ?? {})}
        onChange={(pairs) => {
          onConfigChange({
            ...config,
            headers: Object.fromEntries(pairs),
          })
        }}
      />
    </div>
  )
}

export default ForwarderEditorHttp
