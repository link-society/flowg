import React from 'react'

import InputAdornment from '@mui/material/InputAdornment'
import TextField from '@mui/material/TextField'

import ForwarderConfigModel, {
  ForwarderConfigTypeIconMap,
  ForwarderConfigTypeLabelMap,
} from '@/lib/models/ForwarderConfigModel'
import ForwarderModel from '@/lib/models/ForwarderModel'

import ForwarderEditorAmqp from '@/components/ForwarderEditorAmqp'
import ForwarderEditorClickhouse from '@/components/ForwarderEditorClickhouse'
import ForwarderEditorDatadog from '@/components/ForwarderEditorDatadog'
import ForwarderEditorElastic from '@/components/ForwarderEditorElastic'
import ForwarderEditorHttp from '@/components/ForwarderEditorHttp'
import ForwarderEditorOtlp from '@/components/ForwarderEditorOtlp'
import ForwarderEditorSplunk from '@/components/ForwarderEditorSplunk'
import ForwarderEditorSyslog from '@/components/ForwarderEditorSyslog'

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
