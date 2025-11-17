import Divider from '@mui/material/Divider'
import InputAdornment from '@mui/material/InputAdornment'
import TextField from '@mui/material/TextField'

import ForwarderConfigAmqpModel from '@/lib/models/ForwarderConfigAmqpModel'
import { ForwarderConfigTypeLabelMap } from '@/lib/models/ForwarderConfigModel'

import ForwarderIconAmqp from '@/components/ForwarderIconAmqp'

type ForwarderEditorAmqpProps = {
  config: ForwarderConfigAmqpModel
  onConfigChange: (config: ForwarderConfigAmqpModel) => void
}

const ForwarderEditorAmqp = ({
  config,
  onConfigChange,
}: ForwarderEditorAmqpProps) => {
  return (
    <div
      id="container:editor.forwarders.amqp"
      className="flex flex-col items-stretch gap-3"
    >
      <div className="mb-6 shadow">
        <TextField
          label="Forwarder Type"
          variant="outlined"
          className="w-full"
          type="text"
          value={ForwarderConfigTypeLabelMap.amqp}
          disabled
          slotProps={{
            input: {
              startAdornment: (
                <InputAdornment position="start">
                  <ForwarderIconAmqp />
                </InputAdornment>
              ),
            },
          }}
        />
      </div>

      <TextField
        id="input:editor.forwarders.amqp.url"
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
        id="input:editor.forwarders.amqp.exchange"
        label="Exchange"
        variant="outlined"
        type="text"
        value={config.exchange}
        onChange={(e) => {
          onConfigChange({
            ...config,
            exchange: e.target.value,
          })
        }}
      />

      <TextField
        id="input:editor.forwarders.amqp.routing_key"
        label="Routing Key"
        variant="outlined"
        type="text"
        value={config.routing_key}
        onChange={(e) => {
          onConfigChange({
            ...config,
            routing_key: e.target.value,
          })
        }}
      />
    </div>
  )
}

export default ForwarderEditorAmqp
