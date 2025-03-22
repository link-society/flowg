import { Handle, Position, Node, NodeProps, NodeToolbar } from '@xyflow/react'

import ForwardToInboxIcon from '@mui/icons-material/ForwardToInbox'

import TextField from '@mui/material/TextField'

import { OpenForwarderDialog } from '@/components/editors/forwarder/dialog'
import { DeleteNodeButton } from '../delete-btn'

type ForwarderNodeData = Node<{
  forwarder: string
}>

export const ForwarderNode = ({ id, data, selected }: NodeProps<ForwarderNodeData>) => (
  <>
    {selected && (
      <NodeToolbar className="flex flex-row items-center gap-2">
        <OpenForwarderDialog forwarderName={data.forwarder} />
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
        border-4 border-green-900
        shadow-md hover:shadow-lg
        transition-shadow duration-150 ease-in-out
      "
      style={{
        width: '270px',
        height: '100px',
      }}
    >
      <div className="bg-green-700 text-white p-3 flex flex-row items-center">
        <ForwardToInboxIcon />
      </div>
      <div className="p-3 flex flex-row items-center nodrag">
        <TextField
          label="Forwarder"
          type="text"
          value={data.forwarder}
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
