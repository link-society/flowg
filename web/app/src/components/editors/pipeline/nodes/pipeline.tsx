import { Handle, Position, Node, NodeProps, NodeToolbar } from '@xyflow/react'

import AccountTreeIcon from '@mui/icons-material/AccountTree'

import TextField from '@mui/material/TextField'

import { DeleteNodeButton } from '../delete-btn'

type PipelineNodeData = Node<{
  pipeline: string
}>

export const PipelineNode = ({ id, data, selected }: NodeProps<PipelineNodeData>) => (
  <>
    {selected && (
      <NodeToolbar className="flex flex-row items-center gap-2">
        <DeleteNodeButton nodeId={id} />
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
