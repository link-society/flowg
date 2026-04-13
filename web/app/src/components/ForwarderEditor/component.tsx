import React from 'react'

import InputAdornment from '@mui/material/InputAdornment'

import ForwarderConfigModel, {
  ForwarderConfigTypeIconMap,
  ForwarderConfigTypeLabelMap,
} from '@/lib/models/ForwarderConfigModel'

import ForwarderEditorAmqp from '@/components/ForwarderEditorAmqp/component'
import ForwarderEditorClickhouse from '@/components/ForwarderEditorClickhouse/component'
import ForwarderEditorDatadog from '@/components/ForwarderEditorDatadog/component'
import ForwarderEditorElastic from '@/components/ForwarderEditorElastic/component'
import ForwarderEditorHttp from '@/components/ForwarderEditorHttp/component'
import ForwarderEditorOtlp from '@/components/ForwarderEditorOtlp/component'
import ForwarderEditorSplunk from '@/components/ForwarderEditorSplunk/component'
import ForwarderEditorSyslog from '@/components/ForwarderEditorSyslog/component'

import { ForwarderEditorTextField, ForwarderEditorTypeField } from './styles'
import { ForwarderEditorProps } from './types'

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
        <ForwarderEditorTypeField>
          <ForwarderEditorTextField
            label="Forwarder Type"
            variant="outlined"
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
        </ForwarderEditorTypeField>
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
