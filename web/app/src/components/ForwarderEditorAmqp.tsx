import { useEffect } from 'react'

import Divider from '@mui/material/Divider'
import TextField from '@mui/material/TextField'

import { useInput } from '@/lib/hooks/input'

import ForwarderConfigAmqpModel from '@/lib/models/ForwarderConfigAmqpModel'

import * as validators from '@/lib/validators'

type ForwarderEditorAmqpProps = {
  config: ForwarderConfigAmqpModel
  onConfigChange: (config: ForwarderConfigAmqpModel) => void
  onValidationChange: (valid: boolean) => void
}

const ForwarderEditorAmqp = ({
  config,
  onConfigChange,
  onValidationChange,
}: ForwarderEditorAmqpProps) => {
  const [url, setUrl] = useInput(config.url, [
    validators.minLength(1),
    validators.formatUri,
  ])
  const [exchange, setExchange] = useInput(config.exchange, [
    validators.minLength(1),
  ])
  const [routingKey, setRoutingKey] = useInput(config.routing_key)

  useEffect(() => {
    const valid = url.valid && exchange.valid && routingKey.valid
    onValidationChange(valid)

    if (valid) {
      onConfigChange({
        type: 'amqp',
        url: url.value,
        exchange: exchange.value,
        routing_key: routingKey.value,
      })
    }
  }, [url, exchange, routingKey, onValidationChange, onConfigChange])

  return (
    <div
      id="container:editor.forwarders.amqp"
      className="flex flex-col items-stretch gap-3"
    >
      <TextField
        id="input:editor.forwarders.amqp.url"
        label="URL"
        variant="outlined"
        type="text"
        error={!url.valid}
        value={url.value}
        onChange={(e) => {
          setUrl(e.target.value)
        }}
      />

      <Divider />

      <TextField
        id="input:editor.forwarders.amqp.exchange"
        label="Exchange"
        variant="outlined"
        type="text"
        error={!exchange.valid}
        value={exchange.value}
        onChange={(e) => {
          setExchange(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.amqp.routing_key"
        label="Routing Key"
        variant="outlined"
        type="text"
        error={!routingKey.valid}
        value={routingKey.value}
        onChange={(e) => {
          setRoutingKey(e.target.value)
        }}
      />
    </div>
  )
}

export default ForwarderEditorAmqp
