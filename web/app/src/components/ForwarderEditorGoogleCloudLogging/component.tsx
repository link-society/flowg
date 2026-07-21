import { useColorMode } from '@/theme'

import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import Checkbox from '@mui/material/Checkbox'
import FormControlLabel from '@mui/material/FormControlLabel'
import FormGroup from '@mui/material/FormGroup'
import TextField from '@mui/material/TextField'

import Editor, { useMonaco } from '@monaco-editor/react'

import { useInput } from '@/lib/hooks/input'

import {
  vrlThemeDarkDefinition,
  vrlThemeDefinition,
} from '@/lib/vrl-highlighter.ts'

import {ForwarderEditorGoogleCloudLoggingRoot} from './styles'
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
  const [disable_tls, setDisableTls] = useInput(config.disable_tls)
  const [disable_auth, setDisableAuth] = useInput(config.disable_auth)
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
        disable_auth: disable_auth.value,
        disable_tls: disable_tls.value,
        auth_json,
      })
    }
  }, [
    endpoint,
    project_id,
    log_id,
    disable_auth,
    disable_tls,
    auth_json,
  ])

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

      <FormGroup>
        <FormControlLabel
          control={
            <Checkbox
              id="input:editor.forwarders.googlelog.disable_tls"
              checked={disable_tls.value}
              onChange={(e) => {
                setDisableTls(e.target.checked)
              }}
            />
          }
          label={t(
            'components.forwarderEditorGoogleCloudLogging.allowInsecureLabel'
          )}
        />

        <FormControlLabel
          control={
            <Checkbox
              id="input:editor.forwarders.googlelog.disable_auth"
              checked={disable_auth.value}
              onChange={(e) => {
                setDisableAuth(e.target.checked)
              }}
            />
          }
          label={t(
            'components.forwarderEditorGoogleCloudLogging.disableAuthLabel'
          )}
        />
      </FormGroup>

      {!disable_auth.value && (
        <>
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
        </>
      )}
    </ForwarderEditorGoogleCloudLoggingRoot>
  )
}

export default ForwarderEditorGoogleCloudLogging
