import { useEffect } from 'react'
import { useTranslation } from 'react-i18next'

import Checkbox from '@mui/material/Checkbox'
import FormControlLabel from '@mui/material/FormControlLabel'
import MenuItem from '@mui/material/MenuItem'
import Select from '@mui/material/Select'
import TextField from '@mui/material/TextField'

import { useInput } from '@/lib/hooks/input'

import { DynamicField } from '@/lib/models/DynamicField.ts'
import {
  JetStreamAuth,
  JetStreamAuthValues,
} from '@/lib/models/ForwarderConfigJetStreamModel.ts'

import * as validators from '@/lib/validators'

import DynamicFieldControl from '@/components/DynamicFieldControl/component.tsx'
import InputKeyValue from '@/components/InputKeyValue/component.tsx'
import InputList from '@/components/InputList/component.tsx'

import { ForwarderEditorJetStreamRoot } from './styles'
import { ForwarderEditorJetStreamProps } from './types'

const ForwarderEditorJetStream = ({
  config,
  onConfigChange,
  onValidationChange,
}: ForwarderEditorJetStreamProps) => {
  const { t } = useTranslation()

  const [servers, setServers] = useInput(config.servers, [
    validators.items([validators.minLength(1), validators.formatUri]),
  ])

  const [subject, setSubject] = useInput<DynamicField<string>>(config.subject, [
    validators.dynamicField([]),
  ])

  const [body, setBody] = useInput<DynamicField<string>>(config.body, [
    validators.dynamicField([]),
  ])

  const [messageID, setMessageID] = useInput<DynamicField<string>>(
    config.message_id,
    [validators.dynamicField([])]
  )

  const [expectedStream, setExpectedStream] = useInput(
    config.expected_stream,
    []
  )
  const [publishTimeout, setPublishTimeout] = useInput(
    config.publish_timeout,
    []
  )

  const [headers, setHeaders] = useInput<Record<string, string>>(
    config.headers ?? {},
    []
  )

  const [authMode, setAuthMode] = useInput<JetStreamAuth>(config.auth_mode, [])
  const [username, setUsername] = useInput(config.username, [])
  const [password, setPassword] = useInput(config.password, [])
  const [token, setToken] = useInput(config.token, [])
  const [credentials, setCredentials] = useInput(config.credentials, [])

  const [tls_ca, setTlsCA] = useInput(config.tls_ca, [])
  const [tls_certificate, setTlsCertificate] = useInput(
    config.tls_certificate,
    []
  )
  const [tls_private_key, setTlsPrivateKey] = useInput(
    config.tls_private_key,
    []
  )
  const [tls_server_name, setTlsServerName] = useInput(
    config.tls_server_name,
    []
  )
  const [tls_insecure_skip_verify, setTlsInsecureSkipVerify] = useInput(
    config.tls_insecure_skip_verify,
    []
  )

  useEffect(() => {
    const valid = true
    onValidationChange(valid)

    if (valid) {
      onConfigChange({
        type: 'jetstream',
        servers: servers.value,
        subject: subject.value,
        body: body.value,
        message_id: messageID.value,
        expected_stream: expectedStream.value,
        publish_timeout: publishTimeout.value,
        headers: headers.value,
        auth_mode: authMode.value,
        username: username.value,
        password: password.value,
        token: token.value,
        credentials: credentials.value,
        tls_ca: tls_ca.value,
        tls_certificate: tls_certificate.value,
        tls_private_key: tls_private_key.value,
        tls_server_name: tls_server_name.value,
        tls_insecure_skip_verify: tls_insecure_skip_verify.value,
      })
    }
  }, [
    servers,
    subject,
    body,
    messageID,
    expectedStream,
    publishTimeout,
    headers,
    authMode,
    username,
    password,
    token,
    credentials,
    tls_ca,
    tls_certificate,
    tls_private_key,
    tls_server_name,
    tls_insecure_skip_verify,
  ])

  return (
    <ForwarderEditorJetStreamRoot id="container:editor.forwarders.azuremonitor">
      <DynamicFieldControl
        id="input:editor.forwarders.jetstream.subject"
        label={t('components.forwarderEditorJetStream.subjectLabel')}
        multiline
        variant="outlined"
        error={!subject.valid}
        value={subject.value}
        onChange={setSubject}
      />

      <InputList
        id="editor.forwarders.jetstream.servers"
        itemLabel={t('components.forwarderEditorJetStream.serverItemLabel')}
        items={servers.value}
        itemValidators={[validators.minLength(1), validators.formatUri]}
        onChange={setServers}
      />

      <DynamicFieldControl
        id="input:editor.forwarders.jetstream.body"
        label={t('components.forwarderEditorJetStream.bodyLabel')}
        multiline
        variant="outlined"
        error={!body.valid}
        value={body.value}
        onChange={setBody}
      />

      <DynamicFieldControl
        id="input:editor.forwarders.jetstream.message_id"
        label={t('components.forwarderEditorJetStream.messageIDLabel')}
        multiline
        variant="outlined"
        error={!messageID.valid}
        value={messageID.value}
        onChange={setMessageID}
      />

      <TextField
        id="input:editor.forwarders.jetstream.expected_stream"
        label={t('components.forwarderEditorJetStream.expectedStreamLabel')}
        variant="outlined"
        error={!expectedStream.valid}
        value={expectedStream.value}
        onChange={(e) => {
          setExpectedStream(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.jetstream.publish_timeout"
        label={t('components.forwarderEditorJetStream.PublishTimeoutLabel')}
        variant="outlined"
        type="number"
        error={!publishTimeout.valid}
        value={publishTimeout.value}
        onChange={(e) => {
          setPublishTimeout(+e.target.value)
        }}
      />

      <InputKeyValue
        id="input:editor.forwarders.jetstream.headers"
        keyLabel={t('components.forwarderEditorJetStream.headerNameLabel')}
        valueLabel={t('components.forwarderEditorJetStream.headerValueLabel')}
        keyValues={Object.entries(headers.value ?? {})}
        onChange={(pairs) => {
          setHeaders(Object.fromEntries(pairs))
        }}
      />

      <Select<JetStreamAuth>
        labelId="label:editor.forwarders.jetstream.auth_mode"
        id="select:editor.forwarders.jetstream.auth_mode"
        value={authMode.value}
        label={t('components.forwarderEditorJetStream.authModeLabel')}
        onChange={(e) => {
          setAuthMode(e.target.value as JetStreamAuth)
        }}
      >
        {JetStreamAuthValues.map((t) => (
          <MenuItem
            id={`option:editor.forwarders.jetstream.auth_mode.${t}`}
            key={t}
            value={t}
          >
            {t.toUpperCase()}
          </MenuItem>
        ))}
      </Select>

      <TextField
        id="input:editor.forwarders.jetstream.username"
        label={t('components.forwarderEditorJetStream.usernameLabel')}
        variant="outlined"
        error={!username.valid}
        value={username.value}
        onChange={(e) => {
          setUsername(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.jetstream.password"
        label={t('components.forwarderEditorJetStream.passwordLabel')}
        variant="outlined"
        type="password"
        error={!password.valid}
        value={password.value}
        onChange={(e) => {
          setPassword(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.jetstream.token"
        label={t('components.forwarderEditorJetStream.tokenLabel')}
        variant="outlined"
        error={!token.valid}
        value={token.value}
        onChange={(e) => {
          setToken(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.jetstream.credentials"
        label={t('components.forwarderEditorJetStream.credentialsLabel')}
        variant="outlined"
        error={!credentials.valid}
        value={credentials.value}
        onChange={(e) => {
          setCredentials(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.jetstream.tls_ca"
        label={t('components.forwarderEditorJetStream.tlsCALabel')}
        variant="outlined"
        error={!tls_ca.valid}
        value={tls_ca.value}
        onChange={(e) => {
          setTlsCA(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.jetstream.tls_certificate"
        label={t('components.forwarderEditorJetStream.tlsCertificateLabel')}
        variant="outlined"
        error={!tls_certificate.valid}
        value={tls_certificate.value}
        onChange={(e) => {
          setTlsCertificate(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.jetstream.tls_private_key"
        label={t('components.forwarderEditorJetStream.tlsPrivateKeyLabel')}
        variant="outlined"
        error={!tls_private_key.valid}
        value={tls_private_key.value}
        onChange={(e) => {
          setTlsPrivateKey(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.jetstream.tls_server_name"
        label={t('components.forwarderEditorJetStream.tlsServerNameLabel')}
        variant="outlined"
        error={!tls_server_name.valid}
        value={tls_server_name.value}
        onChange={(e) => {
          setTlsServerName(e.target.value)
        }}
      />

      <FormControlLabel
        control={
          <Checkbox
            id="input:editor.forwarders.jetstream.tls_insecure_skip_verify"
            checked={tls_insecure_skip_verify.value}
            onChange={(e) => {
              setTlsInsecureSkipVerify(e.target.checked)
            }}
          />
        }
        label={t(
          'components.forwarderEditorJetStream.tlsInsecureSkipVerifyLabel'
        )}
      />
    </ForwarderEditorJetStreamRoot>
  )
}

export default ForwarderEditorJetStream
