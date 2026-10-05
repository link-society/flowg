import { useTranslation } from 'react-i18next'

import TextField from '@mui/material/TextField'
import { useTheme } from '@mui/material/styles'

import AccountTreeIcon from '@mui/icons-material/AccountTree'

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

import { PipelineNodePipelineData } from './types'

const PipelineNodePipeline = ({
  data,
  selected,
}: NodeProps<PipelineNodePipelineData>) => {
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
      <NodeRoot borderColor={theme.tokens.colors.nodePipelineBorder}>
        <NodeIcon bgColor={theme.tokens.colors.nodePipelineBg}>
          <AccountTreeIcon />
        </NodeIcon>
        <NodeBody className="nodrag">
          <TextField
            label={t('components.pipelineNodePipeline.label')}
            type="text"
            value={data.pipeline}
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

export default PipelineNodePipeline
