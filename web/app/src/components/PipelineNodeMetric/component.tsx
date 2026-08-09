import React, { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import TextField from '@mui/material/TextField'

import BarChartIcon from '@mui/icons-material/BarChart'

import { Handle, NodeProps, Position } from '@xyflow/react'

import { usePipelineEditorHooks } from '@/lib/hooks/pipeline-editor'

import PipelineDeleteNodeButton from '@/components/PipelineDeleteNodeButton/component'
import PipelineTraceNodeButton from '@/components/PipelineTraceNodeButton/component'
import PipelineTraceNodeIndicator from '@/components/PipelineTraceNodeIndicator/component'

import { NodeBody, NodeIcon, NodeRoot, ToolbarRow, handleStyle } from './styles'
import { PipelineNodeMetricData } from './types'

const PipelineNodeMetric = ({
  id,
  data,
  selected,
}: NodeProps<PipelineNodeMetricData>) => {
  const { t } = useTranslation()
  const { setNodes } = usePipelineEditorHooks()

  const [name, setName] = useState(data.name)

  const onChange: React.ChangeEventHandler<HTMLInputElement> = (evt) => {
    setName(evt.target.value)
  }

  useEffect(() => {
    setNodes((prevNodes) => {
      const newNodes = [...prevNodes]

      for (const node of newNodes) {
        if (node.id === id) {
          node.data = { name }
        }
      }

      return newNodes
    })
  }, [id, name])

  return (
    <>
      {selected && (
        <ToolbarRow>
          <PipelineDeleteNodeButton nodeId={id} />
          {data.traces && <PipelineTraceNodeButton traces={data.traces} />}
        </ToolbarRow>
      )}

      <Handle type="target" position={Position.Left} style={handleStyle} />
      <NodeRoot>
        <NodeIcon>
          <BarChartIcon />
        </NodeIcon>
        <NodeBody className="nodrag">
          <TextField
            label={t('components.pipelineNodeMetric.label')}
            type="text"
            value={name}
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

      <PipelineTraceNodeIndicator
        status={
          data.traces
            ? data.traces.some((trace) => trace.error)
              ? 'error'
              : 'success'
            : null
        }
      />
    </>
  )
}

export default PipelineNodeMetric
