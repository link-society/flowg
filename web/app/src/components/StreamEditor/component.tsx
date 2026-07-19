import { useTranslation } from 'react-i18next'

import Divider from '@mui/material/Divider'
import TextField from '@mui/material/TextField'
import Typography from '@mui/material/Typography'

import ListEdit from '@/components/ListEdit/component'

import {
  StreamEditorHint,
  StreamEditorPanel,
  StreamEditorPanelBody,
  StreamEditorPanelHeader,
  StreamEditorProgress,
  StreamEditorRoot,
  StreamEditorUsageRow,
} from './styles'
import { StreamEditorProps } from './types'

const StreamEditor = ({
  streamConfig,
  storageUsage,
  onStreamConfigChange,
}: StreamEditorProps) => {
  const { t } = useTranslation()
  const usageMB = storageUsage / (1024 * 1024)
  const usagePercent = (usageMB * 100) / streamConfig.size

  return (
    <StreamEditorRoot>
      <StreamEditorPanel>
        <StreamEditorPanelHeader>
          <Typography variant="titleSm">
            {t('components.streamEditor.retentionTitle')}
          </Typography>
        </StreamEditorPanelHeader>
        <Divider />
        <StreamEditorPanelBody>
          <StreamEditorUsageRow>
            <Typography variant="text">
              {t('components.streamEditor.usageLabel', {
                usage: usageMB.toFixed(2),
              })}
            </Typography>
            {streamConfig.size > 0 && (
              <StreamEditorProgress
                variant="determinate"
                color={usagePercent < 100 ? 'primary' : 'error'}
                value={Math.round(Math.min(usagePercent, 100) * 100) / 100}
              />
            )}
          </StreamEditorUsageRow>

          <TextField
            id="input:editor.streams.retention-size"
            label={t('components.streamEditor.retentionSizeLabel')}
            variant="outlined"
            type="number"
            value={streamConfig.size}
            onChange={(e) => {
              onStreamConfigChange({
                ...streamConfig,
                size: Number(e.target.value),
              })
            }}
          />

          <TextField
            id="input:editor.streams.retention-ttl"
            label={t('components.streamEditor.retentionTtlLabel')}
            variant="outlined"
            type="number"
            value={streamConfig.ttl}
            onChange={(e) => {
              onStreamConfigChange({
                ...streamConfig,
                ttl: Number(e.target.value),
              })
            }}
          />

          <StreamEditorHint variant="text">
            {t('components.streamEditor.disableHintPrefix')} <code>0</code>{' '}
            {t('components.streamEditor.disableHintSuffix')}
          </StreamEditorHint>
        </StreamEditorPanelBody>
      </StreamEditorPanel>

      <StreamEditorPanel>
        <StreamEditorPanelHeader>
          <Typography variant="titleSm">
            {t('components.streamEditor.indexesTitle')}
          </Typography>
        </StreamEditorPanelHeader>
        <Divider />
        <StreamEditorPanelBody>
          <ListEdit
            id="editor.streams.indexed-field"
            list={streamConfig.indexed_fields}
            setList={(list) =>
              onStreamConfigChange({ ...streamConfig, indexed_fields: list })
            }
          />
        </StreamEditorPanelBody>
      </StreamEditorPanel>
    </StreamEditorRoot>
  )
}

export default StreamEditor
