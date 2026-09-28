import { useTranslation } from 'react-i18next'

import TextField from '@mui/material/TextField'
import { useTheme } from '@mui/material/styles'

import DeviceHubIcon from '@mui/icons-material/DeviceHub'

import { Handle, NodeProps, Position } from '@xyflow/react'

import {
  NodeBody,
  NodeIcon,
  NodeRoot,
  ToolbarRow,
  handleStyle,
} from '@/components/PipelineNodeCard/styles'
import PipelineNodeDragHandle from '@/components/PipelineNodeDragHandle/component'
import PipelineTraceNodeButton from '@/components/PipelineTraceNodeButton/component'
import PipelineTraceNodeIndicator from '@/components/PipelineTraceNodeIndicator/component'

import { PipelineNodeSwitchData } from './types'

const PipelineNodeSwitch = ({
  data,
  selected,
}: NodeProps<PipelineNodeSwitchData>) => {
  const { t } = useTranslation()
  const theme = useTheme()

  return (
    <>
      {selected && data.traces && (
        <ToolbarRow>
          <PipelineTraceNodeButton traces={data.traces} />
        </ToolbarRow>
      )}

      <Handle type="target" position={Position.Left} style={handleStyle} />
      <NodeRoot borderColor={theme.tokens.colors.nodeSwitchBorder}>
        <PipelineNodeDragHandle color={theme.tokens.colors.nodeSwitchBg} />
        <NodeIcon bgColor={theme.tokens.colors.nodeSwitchBg}>
          <DeviceHubIcon />
        </NodeIcon>
        <NodeBody className="nodrag">
          <TextField
            label={t('components.pipelineNodeSwitch.label')}
            type="text"
            value={data.condition}
            slotProps={{
              input: {
                readOnly: true,
                sx: { fontFamily: 'monospace' },
              },
            }}
            variant="outlined"
          />
        </NodeBody>
      </NodeRoot>
      <Handle type="source" position={Position.Right} style={handleStyle} />

      <PipelineTraceNodeIndicator
        status={
          data.traces
            ? data.traces.some((trace) => trace.error)
              ? 'error'
              : 'success'
            : null
        }
      />
    </>
  )
}

export default PipelineNodeSwitch
