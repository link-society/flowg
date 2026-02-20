import TextField from '@mui/material/TextField'

import InputIcon from '@mui/icons-material/Input'

import { Handle, Node, NodeProps, NodeToolbar, Position } from '@xyflow/react'

import { NodeTrace } from '@/lib/models/PipelineTrace.ts'

import PipelineTraceNodeButton from '@/components/PipelineTraceNodeButton.tsx'
import PipelineTraceNodeIndicator from '@/components/PipelineTraceNodeIndicator.tsx'

type PipelineNodeSourceData = Node<{
  type: string
  trace: NodeTrace | null
}>

const PipelineNodeSource = ({
  selected,
  data,
}: NodeProps<PipelineNodeSourceData>) => (
  <>
    {selected && data.trace && (
      <NodeToolbar className="flex flex-row items-center gap-2">
        <PipelineTraceNodeButton trace={data.trace} />
      </NodeToolbar>
    )}

    <div
      className="
        flex flex-row items-stretch gap-2
        bg-white
        border-4 border-orange-700
        shadow-md hover:shadow-lg
        transition-shadow duration-150 ease-in-out
      "
      style={{
        width: '270px',
        height: '100px',
      }}
    >
      <div className="bg-orange-500 text-white p-3 flex flex-row items-center">
        <InputIcon />
      </div>
      <div className="p-3 flex flex-row items-center nodrag">
        <TextField
          label="Source"
          type="text"
          value={data.type.toUpperCase()}
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

    <PipelineTraceNodeIndicator trace={data.trace} />
  </>
)

export default PipelineNodeSource
