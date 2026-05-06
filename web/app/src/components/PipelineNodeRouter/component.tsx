import TextField from '@mui/material/TextField'

import StorageIcon from '@mui/icons-material/Storage'

import { Handle, NodeProps, Position } from '@xyflow/react'

import DialogStreamEditor from '@/components/DialogStreamEditor/component'
import PipelineDeleteNodeButton from '@/components/PipelineDeleteNodeButton/component'
import PipelineTraceNodeButton from '@/components/PipelineTraceNodeButton/component'
import PipelineTraceNodeIndicator from '@/components/PipelineTraceNodeIndicator/component'

import { NodeBody, NodeIcon, NodeRoot, ToolbarRow, handleStyle } from './styles'
import { PipelineNodeRouterData } from './types'

const PipelineNodeRouter = ({
  id,
  data,
  selected,
}: NodeProps<PipelineNodeRouterData>) => (
  <>
    {selected && (
      <ToolbarRow>
        <DialogStreamEditor stream={data.stream} />
        <PipelineDeleteNodeButton nodeId={id} />
        {data.traces && <PipelineTraceNodeButton traces={data.traces} />}
      </ToolbarRow>
    )}

    <Handle type="target" position={Position.Left} style={handleStyle} />
    <NodeRoot>
      <NodeIcon>
        <StorageIcon />
      </NodeIcon>
      <NodeBody className="nodrag">
        <TextField
          label="Stream"
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

export default PipelineNodeRouter
