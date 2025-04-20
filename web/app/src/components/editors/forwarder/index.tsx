import { ForwarderConfigModel, ForwarderModel } from '@/lib/models/forwarder'

import { AmqpForwarderEditor } from './amqp'
import { DatadogForwarderEditor } from './datadog'
import { HttpForwarderEditor } from './http'
import { SyslogForwarderEditor } from './syslog'

type ForwarderEditorProps = {
  forwarder: ForwarderModel
  onForwarderChange: (forwarder: ForwarderModel) => void
}

export const ForwarderEditor = ({
  forwarder,
  onForwarderChange,
}: ForwarderEditorProps) => {
  const onConfigChange = (config: ForwarderConfigModel) => {
    onForwarderChange({
      ...forwarder,
      config,
    })
  }

  switch (forwarder.config.type) {
    case 'http':
      return (
        <HttpForwarderEditor
          config={forwarder.config}
          onConfigChange={onConfigChange}
        />
      )

    case 'syslog':
      return (
        <SyslogForwarderEditor
          config={forwarder.config}
          onConfigChange={onConfigChange}
        />
      )

    case 'datadog':
      return (
        <DatadogForwarderEditor
          config={forwarder.config}
          onConfigChange={onConfigChange}
        />
      )

    case 'amqp':
      return (
        <AmqpForwarderEditor
          config={forwarder.config}
          onConfigChange={onConfigChange}
        />
      )
  }
}
