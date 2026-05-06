import TextField from '@mui/material/TextField'

import InputIcon from '@mui/icons-material/Input'

import { Handle, NodeProps, Position } from '@xyflow/react'

import PipelineTraceNodeButton from '@/components/PipelineTraceNodeButton/component'
import PipelineTraceNodeIndicator from '@/components/PipelineTraceNodeIndicator/component'

import { NodeBody, NodeIcon, NodeRoot, ToolbarRow, handleStyle } from './styles'
import { PipelineNodeSourceData } from './types'

const PipelineNodeSource = ({
  selected,
  data,
}: NodeProps<PipelineNodeSourceData>) => (
  <>
    {selected && data.traces && (
      <ToolbarRow>
        <PipelineTraceNodeButton traces={data.traces} />
      </ToolbarRow>
    )}

    <NodeRoot>
      <NodeIcon>
        <InputIcon />
      </NodeIcon>
      <NodeBody className="nodrag">
        <TextField
          label="Source"
          type="text"
          value={data.type.toUpperCase()}
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

export default PipelineNodeSource
