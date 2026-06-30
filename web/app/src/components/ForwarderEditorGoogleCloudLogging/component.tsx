import { useEffect } from 'react'

import TextField from '@mui/material/TextField'

import { useInput } from '@/lib/hooks/input'

import { ForwarderEditorGoogleCloudLoggingRoot } from './styles'
import { ForwarderEditorGoogleCloudLoggingProps } from './types'

const ForwarderEditorGoogleCloudLogging = ({
  config,
  onConfigChange,
  onValidationChange,
}: ForwarderEditorGoogleCloudLoggingProps) => {
  const [endpoint, setEndpoint] = useInput(config.endpoint)
  const [project_id, setProjectID] = useInput(config.project_id)
  const [log_id, setLogID] = useInput(config.log_id)

  useEffect(() => {
    const valid = true
    onValidationChange(valid)

    if (valid) {
      onConfigChange({
        type: 'googlecloudlogging',
        endpoint: endpoint.value,
        project_id: project_id.value,
        log_id: log_id.value,
      })
    }
  }, [project_id, log_id, endpoint])

  return (
    <ForwarderEditorGoogleCloudLoggingRoot id="container:editor.forwarders.googlelog">
      <TextField
        id="input:editor.forwarders.googlelog.endpoint"
        label="Endpoint"
        variant="outlined"
        type="text"
        error={!endpoint.valid}
        value={endpoint.value}
        onChange={(e) => {
          setEndpoint(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.googlelog.project_id"
        label="Project ID"
        variant="outlined"
        type="text"
        error={!project_id.valid}
        value={project_id.value}
        onChange={(e) => {
          setProjectID(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.googlelog.log_id"
        label="Log ID"
        variant="outlined"
        type="text"
        error={!log_id.valid}
        value={log_id.value}
        onChange={(e) => {
          setLogID(e.target.value)
        }}
      />
    </ForwarderEditorGoogleCloudLoggingRoot>
  )
}

export default ForwarderEditorGoogleCloudLogging
