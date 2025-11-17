import React from 'react'

import ForwarderConfigModel from '@/lib/models/ForwarderConfigModel'
import ForwarderModel from '@/lib/models/ForwarderModel'

import ForwarderEditorAmqp from './ForwarderEditorAmqp'
import ForwarderEditorClickhouse from './ForwarderEditorClickhouse'
import ForwarderEditorDatadog from './ForwarderEditorDatadog'
import ForwarderEditorElastic from './ForwarderEditorElastic'
import ForwarderEditorHttp from './ForwarderEditorHttp'
import ForwarderEditorOtlp from './ForwarderEditorOtlp'
import ForwarderEditorSplunk from './ForwarderEditorSplunk'
import ForwarderEditorSyslog from './ForwarderEditorSyslog'

const editors = {
  amqp: ForwarderEditorAmqp,
  clickhouse: ForwarderEditorClickhouse,
  datadog: ForwarderEditorDatadog,
  elastic: ForwarderEditorElastic,
  http: ForwarderEditorHttp,
  otlp: ForwarderEditorOtlp,
  splunk: ForwarderEditorSplunk,
  syslog: ForwarderEditorSyslog,
}

type ForwarderEditorProps = {
  forwarder: ForwarderModel
  onForwarderChange: (forwarder: ForwarderModel) => void
}

const ForwarderEditor = ({
  forwarder,
  onForwarderChange,
}: ForwarderEditorProps) => {
  const onConfigChange = (config: ForwarderConfigModel) => {
    onForwarderChange({
      ...forwarder,
      config,
    })
  }

  const EditorComponent = editors[forwarder.config.type] as React.FC<{
    config: typeof forwarder.config
    onConfigChange: (config: typeof forwarder.config) => void
  }>
  return (
    <EditorComponent
      config={forwarder.config}
      onConfigChange={onConfigChange}
    />
  )
}

export default ForwarderEditor
