import TextField from '@mui/material/TextField'

import ForwardToInboxIcon from '@mui/icons-material/ForwardToInbox'

import { Handle, Node, NodeProps, NodeToolbar, Position } from '@xyflow/react'

import DialogForwarderEditor from '@/components/DialogForwarderEditor'
import PipelineDeleteNodeButton from '@/components/PipelineDeleteNodeButton'

type PipelineNodeForwarderData = Node<{
  forwarder: string
}>

const PipelineNodeForwarder = ({
  id,
  data,
  selected,
}: NodeProps<PipelineNodeForwarderData>) => (
  <>
    {selected && (
      <NodeToolbar className="flex flex-row items-center gap-2">
        <DialogForwarderEditor forwarderName={data.forwarder} />
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

export default PipelineNodeForwarder
