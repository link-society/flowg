import Divider from '@mui/material/Divider'
import LinearProgress from '@mui/material/LinearProgress'
import TextField from '@mui/material/TextField'
import Typography from '@mui/material/Typography'

import ListEdit from '@/components/ListEdit/component'

import {
  StreamEditorHint,
  StreamEditorPanel,
  StreamEditorPanelBody,
  StreamEditorPanelHeader,
  StreamEditorRoot,
  StreamEditorUsageRow,
} from './styles'
import { StreamEditorProps } from './types'

const StreamEditor = ({
  streamConfig,
  storageUsage,
  onStreamConfigChange,
}: StreamEditorProps) => {
  const usageMB = storageUsage / (1024 * 1024)
  const usagePercent = (usageMB * 100) / streamConfig.size

  return (
    <StreamEditorRoot>
      <StreamEditorPanel>
        <StreamEditorPanelHeader>
          <Typography variant="titleMd">Retention</Typography>
        </StreamEditorPanelHeader>
        <Divider />
        <StreamEditorPanelBody>
          <StreamEditorUsageRow>
            <Typography variant="text">
              Estimated storage usage: {usageMB.toFixed(2)}MB
            </Typography>
            {streamConfig.size > 0 && (
              <LinearProgress
                sx={{ flexGrow: 1, height: '20px' }}
                variant="determinate"
                color={usagePercent < 100 ? 'primary' : 'error'}
                value={Math.round(Math.min(usagePercent, 100) * 100) / 100}
              />
            )}
          </StreamEditorUsageRow>

          <TextField
            id="input:editor.streams.retention-size"
            label="Retention size (in MB)"
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
            label="Retention time (in seconds)"
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
            Use <code>0</code> to disable
          </StreamEditorHint>
        </StreamEditorPanelBody>
      </StreamEditorPanel>

      <StreamEditorPanel>
        <StreamEditorPanelHeader>
          <Typography variant="titleMd">Indexes</Typography>
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
