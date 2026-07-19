import { useEffect } from 'react'
import { useTranslation } from 'react-i18next'

import Divider from '@mui/material/Divider'
import TextField from '@mui/material/TextField'

import { useInput } from '@/lib/hooks/input'

import * as validators from '@/lib/validators'

import { ForwarderEditorAwsCloudWatchRoot } from './styles'
import { ForwarderEditorAwsCloudWatchProps } from './types'

const ForwarderEditorAwsCloudWatch = ({
  config,
  onConfigChange,
  onValidationChange,
}: ForwarderEditorAwsCloudWatchProps) => {
  const { t } = useTranslation()
  const [app_id, setAppID] = useInput(config.app_id)

  const [endpoint, setEndpoint] = useInput(config.endpoint, [
    validators.minLength(1),
    validators.formatUri,
  ])

  const [region, setRegion] = useInput(config.region)
  const [akid, setAkid] = useInput(config.access_key_id)
  const [access_key, setAccessKey] = useInput(config.secret_access_key)
  const [token, setToken] = useInput(config.session_token)
  const [group, setGroup] = useInput(config.group, [validators.minLength(1)])
  const [stream, setStream] = useInput(config.stream, [validators.minLength(1)])

  useEffect(() => {
    const valid = true
    onValidationChange(valid)

    if (valid) {
      onConfigChange({
        type: 'awscloudwatch',
        app_id: app_id.value,
        endpoint: endpoint.value,
        region: region.value,
        access_key_id: akid.value,
        secret_access_key: access_key.value,
        session_token: token.value,
        group: group.value,
        stream: stream.value,
      })
    }
  }, [app_id, endpoint, region, akid, access_key, token, group, stream])

  return (
    <ForwarderEditorAwsCloudWatchRoot id="container:editor.forwarders.cloudwatch">
      <TextField
        id="input:editor.forwarders.cloudwatch.app_id"
        label={t('components.forwarderEditorAwsCloudWatch.appIdLabel')}
        variant="outlined"
        type="text"
        error={!app_id.valid}
        value={app_id.value}
        onChange={(e) => {
          setAppID(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.cloudwatch.webhook_url"
        label={t('components.forwarderEditorAwsCloudWatch.endpointLabel')}
        variant="outlined"
        type="text"
        error={!endpoint.valid}
        value={endpoint.value}
        onChange={(e) => {
          setEndpoint(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.cloudwatch.region"
        label={t('components.forwarderEditorAwsCloudWatch.regionLabel')}
        variant="outlined"
        type="text"
        error={!region.valid}
        value={region.value}
        onChange={(e) => {
          setRegion(e.target.value)
        }}
      />

      <Divider />

      <TextField
        id="input:editor.forwarders.cloudwatch.akid"
        label={t('components.forwarderEditorAwsCloudWatch.akidLabel')}
        variant="outlined"
        type="text"
        error={!akid.valid}
        value={akid.value}
        onChange={(e) => {
          setAkid(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.cloudwatch.access_key"
        label={t('components.forwarderEditorAwsCloudWatch.accessKeyLabel')}
        variant="outlined"
        type="text"
        error={!access_key.valid}
        value={access_key.value}
        onChange={(e) => {
          setAccessKey(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.cloudwatch.token"
        label={t('components.forwarderEditorAwsCloudWatch.tokenLabel')}
        variant="outlined"
        type="text"
        error={!token.valid}
        value={token.value}
        onChange={(e) => {
          setToken(e.target.value)
        }}
      />

      <Divider />

      <TextField
        id="input:editor.forwarders.cloudwatch.group"
        label={t('components.forwarderEditorAwsCloudWatch.groupLabel')}
        variant="outlined"
        type="text"
        error={!group.valid}
        value={group.value}
        onChange={(e) => {
          setGroup(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.cloudwatch.stream"
        label={t('components.forwarderEditorAwsCloudWatch.streamLabel')}
        variant="outlined"
        type="text"
        error={!stream.valid}
        value={stream.value}
        onChange={(e) => {
          setStream(e.target.value)
        }}
      />
    </ForwarderEditorAwsCloudWatchRoot>
  )
}

export default ForwarderEditorAwsCloudWatch
