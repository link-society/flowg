import TextField from '@mui/material/TextField'

import AccountTreeIcon from '@mui/icons-material/AccountTree'

import { Handle, NodeProps, Position } from '@xyflow/react'

import PipelineDeleteNodeButton from '@/components/PipelineDeleteNodeButton/component'
import PipelineTraceNodeButton from '@/components/PipelineTraceNodeButton/component'
import PipelineTraceNodeIndicator from '@/components/PipelineTraceNodeIndicator/component'

import { NodeBody, NodeIcon, NodeRoot, ToolbarRow, handleStyle } from './styles'
import { PipelineNodePipelineData } from './types'

const PipelineNodePipeline = ({
  id,
  data,
  selected,
}: NodeProps<PipelineNodePipelineData>) => (
  <>
    {selected && (
      <ToolbarRow>
        <PipelineDeleteNodeButton nodeId={id} />
        {data.trace && <PipelineTraceNodeButton trace={data.trace} />}
      </ToolbarRow>
    )}
    <Handle type="target" position={Position.Left} style={handleStyle} />
    <NodeRoot>
      <NodeIcon>
        <AccountTreeIcon />
      </NodeIcon>
      <NodeBody className="nodrag">
        <TextField
          label="Pipeline"
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

    <PipelineTraceNodeIndicator trace={data.trace} />
  </>
)

export default PipelineNodePipeline
