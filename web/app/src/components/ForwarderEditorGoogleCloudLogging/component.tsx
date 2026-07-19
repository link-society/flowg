import { useColorMode } from '@/theme'

import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import TextField from '@mui/material/TextField'

import Editor, { useMonaco } from '@monaco-editor/react'

import { useInput } from '@/lib/hooks/input'

import {
  vrlThemeDarkDefinition,
  vrlThemeDefinition,
} from '@/lib/vrl-highlighter.ts'

import { ForwarderEditorGoogleCloudLoggingRoot } from './styles'
import { ForwarderEditorGoogleCloudLoggingProps } from './types'

const ForwarderEditorGoogleCloudLogging = ({
  config,
  onConfigChange,
  onValidationChange,
}: ForwarderEditorGoogleCloudLoggingProps) => {
  const { t } = useTranslation()
  const [endpoint, setEndpoint] = useInput(config.endpoint)
  const [project_id, setProjectID] = useInput(config.project_id)
  const [log_id, setLogID] = useInput(config.log_id)
  const [auth_json, setAuthJson] = useState(config.auth_json)

  const monaco = useMonaco()
  const { mode } = useColorMode()

  useEffect(() => {
    if (!monaco) return

    monaco.languages.register({ id: 'json' })
    monaco.editor.defineTheme('vrl-theme-light', vrlThemeDefinition as any)
    monaco.editor.defineTheme('vrl-theme-dark', vrlThemeDarkDefinition as any)
    monaco.editor.setTheme(
      mode === 'dark' ? 'vrl-theme-dark' : 'vrl-theme-light'
    )
  }, [monaco, mode])

  useEffect(() => {
    const valid = true
    onValidationChange(valid)

    if (valid) {
      onConfigChange({
        type: 'googlecloudlogging',
        endpoint: endpoint.value,
        project_id: project_id.value,
        log_id: log_id.value,
        auth_json,
      })
    }
  }, [endpoint, project_id, log_id, auth_json])

  return (
    <ForwarderEditorGoogleCloudLoggingRoot id="container:editor.forwarders.googlelog">
      <TextField
        id="input:editor.forwarders.googlelog.endpoint"
        label={t('components.forwarderEditorGoogleCloudLogging.endpointLabel')}
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
        label={t('components.forwarderEditorGoogleCloudLogging.projectIdLabel')}
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
        label={t('components.forwarderEditorGoogleCloudLogging.logIdLabel')}
        variant="outlined"
        type="text"
        error={!log_id.valid}
        value={log_id.value}
        onChange={(e) => {
          setLogID(e.target.value)
        }}
      />

      <label>
        {t('components.forwarderEditorGoogleCloudLogging.authJsonLabel')}
      </label>
      <Editor
        defaultValue={auth_json}
        defaultLanguage="json"
        height="10rem"
        theme={mode === 'dark' ? 'vrl-theme-dark' : 'vrl-theme-light'}
        onChange={setAuthJson}
        options={{ minimap: { enabled: false } }}
      />
    </ForwarderEditorGoogleCloudLoggingRoot>
  )
}

export default ForwarderEditorGoogleCloudLogging
