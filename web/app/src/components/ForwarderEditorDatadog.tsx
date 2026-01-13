import { useEffect } from 'react'

import Divider from '@mui/material/Divider'
import TextField from '@mui/material/TextField'

import { useInput } from '@/lib/hooks/input'

import { DynamicField } from '@/lib/models/DynamicField.ts'
import ForwarderConfigDatadogModel from '@/lib/models/ForwarderConfigDatadogModel'

import * as validators from '@/lib/validators'

import DynamicFieldControl from '@/components/DynamicFieldControl.tsx'

type ForwarderEditorDatadogProps = {
  config: ForwarderConfigDatadogModel
  onConfigChange: (config: ForwarderConfigDatadogModel) => void
  onValidationChange: (valid: boolean) => void
}

const ForwarderEditorDatadog = ({
  config,
  onConfigChange,
  onValidationChange,
}: ForwarderEditorDatadogProps) => {
  const [url, setUrl] = useInput(config.url, [
    validators.minLength(1),
    validators.formatUri,
  ])
  const [apiKey, setApiKey] = useInput(config.apiKey, [validators.minLength(1)])

  const [ddsource, setDdsource] = useInput<DynamicField<string>>(
    config.ddsource,
    [validators.dynamicField([])]
  )

  const [ddtags, setDdtags] = useInput<DynamicField<string>>(config.ddtags, [
    validators.dynamicField([]),
  ])

  const [hostname, setHostname] = useInput<DynamicField<string>>(
    config.hostname,
    [validators.dynamicField([])]
  )

  const [message, setMessage] = useInput<DynamicField<string>>(config.message, [
    validators.dynamicField([]),
  ])

  const [service, setService] = useInput<DynamicField<string>>(config.service, [
    validators.dynamicField([]),
  ])

  useEffect(() => {
    const valid = url.valid && apiKey.valid
    onValidationChange(valid)

    if (valid) {
      onConfigChange({
        type: 'datadog',
        url: url.value,
        apiKey: apiKey.value,
        ddsource: ddsource.value,
        ddtags: ddtags.value,
        hostname: hostname.value,
        message: message.value,
        service: service.value,
      })
    }
  }, [url, apiKey])

  return (
    <div
      id="container:editor.forwarders.datadog"
      className="flex flex-col items-stretch gap-3"
    >
      <TextField
        id="input:editor.forwarders.datadog.url"
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
        id="input:editor.forwarders.datadog.apiKey"
        label="ApiKey"
        variant="outlined"
        type="password"
        error={!apiKey.valid}
        value={apiKey.value}
        onChange={(e) => {
          setApiKey(e.target.value)
        }}
      />

      <DynamicFieldControl
        id="input:editor.forwarders.datadog.source"
        label="Source"
        variant="outlined"
        error={!ddsource.valid}
        value={ddsource.value}
        onChange={setDdsource}
      />
      <DynamicFieldControl
        id="input:editor.forwarders.datadog.tags"
        label="Tags"
        variant="outlined"
        error={!ddtags.valid}
        value={ddtags.value}
        onChange={setDdtags}
      />
      <DynamicFieldControl
        id="input:editor.forwarders.datadog.hostname"
        label="Hostname"
        variant="outlined"
        error={!hostname.valid}
        value={hostname.value}
        onChange={setHostname}
      />
      <DynamicFieldControl
        id="input:editor.forwarders.datadog.message"
        label="Message"
        variant="outlined"
        error={!message.valid}
        value={message.value}
        onChange={setMessage}
      />
      <DynamicFieldControl
        id="input:editor.forwarders.datadog.service"
        label="Service"
        variant="outlined"
        error={!service.valid}
        value={service.value}
        onChange={setService}
      />
    </div>
  )
}

export default ForwarderEditorDatadog
