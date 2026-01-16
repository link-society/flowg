import { useEffect } from 'react'

import Divider from '@mui/material/Divider'
import TextField from '@mui/material/TextField'

import { useInput } from '@/lib/hooks/input'

import { DynamicField } from '@/lib/models/DynamicField.ts'
import ForwarderConfigAmqpModel from '@/lib/models/ForwarderConfigAmqpModel'

import * as validators from '@/lib/validators'

import DynamicFieldControl from '@/components/DynamicFieldControl.tsx'

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
  const [exchange, setExchange] = useInput<DynamicField<string>>(
    config.exchange,
    [validators.dynamicField([]), validators.minLength(1)]
  )
  const [routingKey, setRoutingKey] = useInput<DynamicField<string>>(
    config.routing_key,
    [validators.dynamicField([])]
  )

  const [body, setBody] = useInput<DynamicField<string>>(config.body, [
    validators.dynamicField([]),
  ])

  useEffect(() => {
    const valid = url.valid && exchange.valid && routingKey.valid
    onValidationChange(valid)

    if (valid) {
      onConfigChange({
        type: 'amqp',
        url: url.value,
        exchange: exchange.value,
        routing_key: routingKey.value,
        body: body.value,
      })
    }
  }, [url, exchange, routingKey])

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

      <DynamicFieldControl
        id="input:editor.forwarders.amqp.exchange"
        label="Exchange"
        variant="outlined"
        type="text"
        error={!exchange.valid}
        value={exchange.value}
        onChange={setExchange}
      />

      <DynamicFieldControl
        id="input:editor.forwarders.amqp.routing_key"
        label="Routing Key"
        variant="outlined"
        type="text"
        error={!routingKey.valid}
        value={routingKey.value}
        onChange={setRoutingKey}
      />

      <DynamicFieldControl
        id="input:editor.forwarders.amqp.body"
        label="Routing Key"
        multiline
        variant="outlined"
        type="text"
        error={!routingKey.valid}
        value={routingKey.value}
        onChange={setBody}
      />
    </div>
  )
}

export default ForwarderEditorAmqp
