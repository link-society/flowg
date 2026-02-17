import TextField from '@mui/material/TextField'

import AccountTreeIcon from '@mui/icons-material/AccountTree'

import { Handle, Node, NodeProps, NodeToolbar, Position } from '@xyflow/react'

import { NodeTrace } from '@/lib/models/PipelineTrace.ts'

import PipelineDeleteNodeButton from '@/components/PipelineDeleteNodeButton'
import PipelineTraceNodeButton from '@/components/PipelineTraceNodeButton.tsx'

type PipelineNodePipelineData = Node<{
  pipeline: string
  trace: NodeTrace | null
}>

const PipelineNodePipeline = ({
  id,
  data,
  selected,
}: NodeProps<PipelineNodePipelineData>) => (
  <>
    {selected && (
      <NodeToolbar className="flex flex-row items-center gap-2">
        <PipelineDeleteNodeButton nodeId={id} />
        {data.trace && <PipelineTraceNodeButton trace={data.trace} />}
      </NodeToolbar>
    )}
    <Handle
      type="target"
      position={Position.Left}
      style={{
        width: '12px',
        height: '12px',
      }}
    />
    <div
      className="
        flex flex-row items-stretch gap-2
        bg-white
        border-4 border-yellow-600
        shadow-md hover:shadow-lg
        transition-shadow duration-150 ease-in-out
      "
      style={{
        width: '270px',
        height: '100px',
      }}
    >
      <div className="bg-yellow-500 text-white p-3 flex flex-row items-center">
        <AccountTreeIcon />
      </div>
      <div className="p-3 flex flex-row items-center nodrag">
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
      </div>
    </div>

    {data.trace && (
      <div
        style={{
          width: '18px',
          height: '18px',
          position: 'absolute',
          right: '-9px',
          top: '-9px',
          backgroundColor: '#ff4444',
          borderRadius: '50%',
          boxShadow: '-2px 2px 2px #00000055',
        }}
      />
    )}
  </>
)

export default PipelineNodePipeline
