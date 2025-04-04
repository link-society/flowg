import TextField from '@mui/material/TextField'

import StorageIcon from '@mui/icons-material/Storage'

import { Handle, Node, NodeProps, NodeToolbar, Position } from '@xyflow/react'

import { OpenStreamDialog } from '@/components/editors/stream/dialog'

import { DeleteNodeButton } from '../delete-btn'

type RouterNodeData = Node<{
  stream: string
}>

export const RouterNode = ({
  id,
  data,
  selected,
}: NodeProps<RouterNodeData>) => (
  <>
    {selected && (
      <NodeToolbar className="flex flex-row items-center gap-2">
        <OpenStreamDialog stream={data.stream} />
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
        border-4 border-purple-900
        shadow-md hover:shadow-lg
        transition-shadow duration-150 ease-in-out
      "
      style={{
        width: '270px',
        height: '100px',
      }}
    >
      <div className="bg-purple-700 text-white p-3 flex flex-row items-center">
        <StorageIcon />
      </div>
      <div className="p-3 flex flex-row items-center nodrag">
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
      </div>
    </div>
  </>
)
