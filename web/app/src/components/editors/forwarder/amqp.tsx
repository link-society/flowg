import Divider from '@mui/material/Divider'
import InputAdornment from '@mui/material/InputAdornment'
import TextField from '@mui/material/TextField'

import { AmqpIcon } from '@/components/icons/amqp'

import { ForwarderTypeLabelMap } from '@/lib/models/forwarder'
import { AmqpForwarderModel } from '@/lib/models/forwarder/amqp'

type AmqpForwarderEditorProps = {
  config: AmqpForwarderModel
  onConfigChange: (config: AmqpForwarderModel) => void
}

export const AmqpForwarderEditor = ({
  config,
  onConfigChange,
}: AmqpForwarderEditorProps) => {
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
          value={ForwarderTypeLabelMap.amqp}
          disabled
          slotProps={{
            input: {
              startAdornment: (
                <InputAdornment position="start">
                  <AmqpIcon />
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
