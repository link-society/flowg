import React from 'react'

import TextField from '@mui/material/TextField'
import InputAdornment from '@mui/material/InputAdornment'

import ForwarderConfigModel, {
  ForwarderConfigTypeIconMap,
  ForwarderConfigTypeLabelMap,
} from '@/lib/models/ForwarderConfigModel'
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
  onValidationChange: (valid: boolean) => void
  showType?: boolean
}

const ForwarderEditor = ({
  forwarder,
  onForwarderChange,
  onValidationChange,
  showType = true,
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
    onValidationChange: (isValid: boolean) => void
    showType: boolean
  }>

  const typeLabel = ForwarderConfigTypeLabelMap[forwarder.config.type]
  const ForwarderIcon = ForwarderConfigTypeIconMap[forwarder.config.type]

  return (
    <>
      {showType && (
        <div className="mb-6 shadow">
          <TextField
            label="Forwarder Type"
            variant="outlined"
            className="w-full"
            type="text"
            value={typeLabel}
            disabled
            slotProps={{
              input: {
                startAdornment: (
                  <InputAdornment position="start">
                    <ForwarderIcon />
                  </InputAdornment>
                ),
              },
            }}
          />
        </div>
      )}

      <EditorComponent
        config={forwarder.config}
        onConfigChange={onConfigChange}
        onValidationChange={onValidationChange}
        showType={showType}
      />
    </>
  )
}

export default ForwarderEditor
