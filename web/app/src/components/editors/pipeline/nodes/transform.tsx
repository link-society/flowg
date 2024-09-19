import { useProfile } from '@/lib/context/profile'

import { Handle, Position, Node, NodeProps, NodeToolbar } from '@xyflow/react'

import FilterAltIcon from '@mui/icons-material/FilterAlt'

import TextField from '@mui/material/TextField'

import { OpenTransformerDialog } from '@/components/editors/transformer/dialog'

type TransformNodeData = Node<{
  transformer: string
}>

export const TransformNode = ({ data, selected }: NodeProps<TransformNodeData>) => {
  const { permissions } = useProfile()

  return (
    <>
      {selected && permissions.can_edit_transformers && (
        <NodeToolbar>
          <OpenTransformerDialog transformer={data.transformer} />
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
          border-4 border-blue-900
          shadow-md hover:shadow-lg
          transition-shadow duration-150 ease-in-out
        "
        style={{
          width: '270px',
          height: '100px',
        }}
      >
        <div className="bg-blue-700 text-white p-3 flex flex-row items-center">
          <FilterAltIcon />
        </div>
        <div className="p-3 flex flex-row items-center nodrag">
          <TextField
            label="Transformer"
            type="text"
            value={data.transformer}
            slotProps={{
              input: {
                readOnly: true,
              },
            }}
            variant="outlined"
          />
        </div>
      </div>
      <Handle
        type="source"
        position={Position.Right}
        style={{
          width: '12px',
          height: '12px',
        }}
      />
    </>
  )
}
