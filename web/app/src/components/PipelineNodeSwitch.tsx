import React, { useEffect, useState } from 'react'

import TextField from '@mui/material/TextField'

import DeviceHubIcon from '@mui/icons-material/DeviceHub'

import { Handle, Node, NodeProps, NodeToolbar, Position } from '@xyflow/react'

import { usePipelineEditorHooks } from '@/lib/hooks/pipeline-editor'

import { NodeTrace } from '@/lib/models/PipelineTrace.ts'

import PipelineDeleteNodeButton from '@/components/PipelineDeleteNodeButton'
import PipelineTraceNodeButton from '@/components/PipelineTraceNodeButton.tsx'

type PipelineNodeSwitchData = Node<{
  condition: string
  trace: NodeTrace | null
}>

const PipelineNodeSwitch = ({
  id,
  data,
  selected,
}: NodeProps<PipelineNodeSwitchData>) => {
  const { setNodes } = usePipelineEditorHooks()

  const [code, setCode] = useState(data.condition)

  const onChange: React.ChangeEventHandler<HTMLInputElement> = (evt) => {
    setCode(evt.target.value)
  }

  useEffect(() => {
    setNodes((prevNodes) => {
      const newNodes = [...prevNodes]

      for (const node of newNodes) {
        if (node.id === id) {
          node.data = { condition: code }
        }
      }

      return newNodes
    })
  }, [id, code])

  return (
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
          border-4 border-red-700
          shadow-md hover:shadow-lg
          transition-shadow duration-150 ease-in-out
        "
        style={{
          width: '270px',
          height: '100px',
        }}
      >
        <div className="bg-red-600 text-white p-3 flex flex-row items-center">
          <DeviceHubIcon />
        </div>
        <div className="p-3 flex flex-row items-center nodrag">
          <TextField
            label="Condition"
            type="text"
            value={code}
            onChange={onChange}
            slotProps={{
              input: {
                className: 'font-mono',
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

export default PipelineNodeSwitch
