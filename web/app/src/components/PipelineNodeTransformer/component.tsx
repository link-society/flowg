import { useTranslation } from 'react-i18next'

import TextField from '@mui/material/TextField'
import { useTheme } from '@mui/material/styles'

import FilterAltIcon from '@mui/icons-material/FilterAlt'

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

import { PipelineNodeTransformerData } from './types'

const PipelineNodeTransformer = ({
  data,
  selected,
}: NodeProps<PipelineNodeTransformerData>) => {
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
      <NodeRoot borderColor={theme.tokens.colors.nodeTransformerBorder}>
        <NodeIcon bgColor={theme.tokens.colors.nodeTransformerBg}>
          <FilterAltIcon />
        </NodeIcon>
        <NodeBody className="nodrag">
          <TextField
            label={t('components.pipelineNodeTransformer.label')}
            type="text"
            value={data.transformer}
            slotProps={{
              input: {
                readOnly: true,
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

export default PipelineNodeTransformer
