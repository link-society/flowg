import TextField from '@mui/material/TextField'

import AccountTreeIcon from '@mui/icons-material/AccountTree'

import { Handle, Node, NodeProps, NodeToolbar, Position } from '@xyflow/react'

import PipelineDeleteNodeButton from '@/components/PipelineDeleteNodeButton'

type PipelineNodePipelineData = Node<{
  pipeline: string
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
  </>
)

export default PipelineNodePipeline
