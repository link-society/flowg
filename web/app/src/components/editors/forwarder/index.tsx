import { ForwarderConfigModel, ForwarderModel } from '@/lib/models/forwarder'

import { AmqpForwarderEditor } from './amqp'
import { ClickhouseForwarderEditor } from './clickhouse'
import { DatadogForwarderEditor } from './datadog'
import { ElasticForwarderEditor } from './elastic'
import { HttpForwarderEditor } from './http'
import { OtlpForwarderEditor } from './otlp'
import { SplunkForwarderEditor } from './splunk'
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

    case 'splunk':
      return (
        <SplunkForwarderEditor
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

    case 'otlp':
      return (
        <OtlpForwarderEditor
          config={forwarder.config}
          onConfigChange={onConfigChange}
        />
      )

    case 'elastic':
      return (
        <ElasticForwarderEditor
          config={forwarder.config}
          onConfigChange={onConfigChange}
        />
      )

    case 'clickhouse':
      return (
        <ClickhouseForwarderEditor
          config={forwarder.config}
          onConfigChange={onConfigChange}
        />
      )
  }
}
