import { useTranslation } from 'react-i18next'

import TextField from '@mui/material/TextField'
import { useTheme } from '@mui/material/styles'

import StorageIcon from '@mui/icons-material/Storage'

import { Handle, NodeProps, Position } from '@xyflow/react'

import {
  NodeBody,
  NodeIcon,
  NodeRoot,
  ToolbarRow,
  handleStyle,
} from '@/components/PipelineNodeCard/styles'
import PipelineTraceNodeButton from '@/components/PipelineTraceNodeButton/component'
import PipelineTraceNodeIndicator from '@/components/PipelineTraceNodeIndicator/component'

import { PipelineNodeRouterData } from './types'

const PipelineNodeRouter = ({
  data,
  selected,
}: NodeProps<PipelineNodeRouterData>) => {
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
      <NodeRoot borderColor={theme.tokens.colors.nodeRouterBorder}>
        <NodeIcon bgColor={theme.tokens.colors.nodeRouterBg}>
          <StorageIcon />
        </NodeIcon>
        <NodeBody className="nodrag">
          <TextField
            label={t('components.pipelineNodeRouter.label')}
            type="text"
            value={data.stream}
            slotProps={{
              input: {
                readOnly: true,
              },
            }}
            variant="outlined"
          />
        </NodeBody>
      </NodeRoot>

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

export default PipelineNodeRouter
