import TextField from '@mui/material/TextField'

import FilterAltIcon from '@mui/icons-material/FilterAlt'

import { Handle, Node, NodeProps, NodeToolbar, Position } from '@xyflow/react'

import { useProfile } from '@/lib/hooks/profile'

import { NodeTrace } from '@/lib/models/PipelineTrace.ts'

import DialogTransformerEditor from '@/components/DialogTransformerEditor'
import PipelineDeleteNodeButton from '@/components/PipelineDeleteNodeButton'
import PipelineTraceNodeButton from '@/components/PipelineTraceNodeButton.tsx'

type PipelineNodeTransformerData = Node<{
  transformer: string
  trace: NodeTrace | null
}>

const PipelineNodeTransformer = ({
  id,
  data,
  selected,
}: NodeProps<PipelineNodeTransformerData>) => {
  const { permissions } = useProfile()

  return (
    <>
      {selected && permissions.can_edit_transformers && (
        <NodeToolbar className="flex flex-row items-center gap-2">
          <DialogTransformerEditor transformer={data.transformer} />
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
}

export default PipelineNodeTransformer
