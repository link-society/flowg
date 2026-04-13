import React, { useEffect, useState } from 'react'

import TextField from '@mui/material/TextField'

import DeviceHubIcon from '@mui/icons-material/DeviceHub'

import { Handle, NodeProps, Position } from '@xyflow/react'

import { usePipelineEditorHooks } from '@/lib/hooks/pipeline-editor'

import PipelineDeleteNodeButton from '@/components/PipelineDeleteNodeButton/component'
import PipelineTraceNodeButton from '@/components/PipelineTraceNodeButton/component'
import PipelineTraceNodeIndicator from '@/components/PipelineTraceNodeIndicator/component'

import { NodeBody, NodeIcon, NodeRoot, ToolbarRow, handleStyle } from './styles'
import { PipelineNodeSwitchData } from './types'

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
        <ToolbarRow>
          <PipelineDeleteNodeButton nodeId={id} />
          {data.trace && <PipelineTraceNodeButton trace={data.trace} />}
        </ToolbarRow>
      )}

      <Handle type="target" position={Position.Left} style={handleStyle} />
      <NodeRoot>
        <NodeIcon>
          <DeviceHubIcon />
        </NodeIcon>
        <NodeBody className="nodrag">
          <TextField
            label="Condition"
            type="text"
            value={code}
            onChange={onChange}
            slotProps={{
              input: {
                sx: { fontFamily: 'monospace' },
              },
            }}
            variant="outlined"
          />
        </NodeBody>
      </NodeRoot>
      <Handle type="source" position={Position.Right} style={handleStyle} />

      <PipelineTraceNodeIndicator trace={data.trace} />
    </>
  )
}

export default PipelineNodeSwitch
