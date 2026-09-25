import { useTranslation } from 'react-i18next'

import TextField from '@mui/material/TextField'

import FilterAltIcon from '@mui/icons-material/FilterAlt'

import { Handle, NodeProps, Position } from '@xyflow/react'

import PipelineTraceNodeButton from '@/components/PipelineTraceNodeButton/component'
import PipelineTraceNodeIndicator from '@/components/PipelineTraceNodeIndicator/component'

import { NodeBody, NodeIcon, NodeRoot, ToolbarRow, handleStyle } from './styles'
import { PipelineNodeTransformerData } from './types'

const PipelineNodeTransformer = ({
  data,
  selected,
}: NodeProps<PipelineNodeTransformerData>) => {
  const { t } = useTranslation()

  return (
    <>
      {selected && data.traces && (
        <ToolbarRow>
          <PipelineTraceNodeButton traces={data.traces} />
        </ToolbarRow>
      )}

      <Handle type="target" position={Position.Left} style={handleStyle} />
      <NodeRoot>
        <NodeIcon>
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
