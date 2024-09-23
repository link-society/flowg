import { Handle, Position, Node, NodeProps, NodeToolbar } from '@xyflow/react'

import NotificationsActiveIcon from '@mui/icons-material/NotificationsActive'

import TextField from '@mui/material/TextField'

import { OpenAlertDialog } from '@/components/editors/alert/dialog'
import { DeleteNodeButton } from '../delete-btn'

type AlertNodeData = Node<{
  alert: string
}>

export const AlertNode = ({ id, data, selected }: NodeProps<AlertNodeData>) => (
  <>
    {selected && (
      <NodeToolbar className="flex flex-row items-center gap-2">
        <OpenAlertDialog alert={data.alert} />
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
        <NotificationsActiveIcon />
      </div>
      <div className="p-3 flex flex-row items-center nodrag">
        <TextField
          label="Alert"
          type="text"
          value={data.alert}
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
