import TextField from '@mui/material/TextField'

import FilterAltIcon from '@mui/icons-material/FilterAlt'

import { Handle, NodeProps, Position } from '@xyflow/react'

import { useProfile } from '@/lib/hooks/profile'

import DialogTransformerEditor from '@/components/DialogTransformerEditor/component'
import PipelineDeleteNodeButton from '@/components/PipelineDeleteNodeButton/component'
import PipelineTraceNodeButton from '@/components/PipelineTraceNodeButton/component'
import PipelineTraceNodeIndicator from '@/components/PipelineTraceNodeIndicator/component'

import { NodeBody, NodeIcon, NodeRoot, ToolbarRow, handleStyle } from './styles'
import { PipelineNodeTransformerData } from './types'

const PipelineNodeTransformer = ({
  id,
  data,
  selected,
}: NodeProps<PipelineNodeTransformerData>) => {
  const { permissions } = useProfile()

  return (
    <>
      {selected && permissions.can_edit_transformers && (
        <ToolbarRow>
          <DialogTransformerEditor transformer={data.transformer} />
          <PipelineDeleteNodeButton nodeId={id} />
          {data.traces && <PipelineTraceNodeButton traces={data.traces} />}
        </ToolbarRow>
      )}

      <Handle type="target" position={Position.Left} style={handleStyle} />
      <NodeRoot>
        <NodeIcon>
          <FilterAltIcon />
        </NodeIcon>
        <NodeBody className="nodrag">
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
        </NodeBody>
      </NodeRoot>
      <Handle type="source" position={Position.Right} style={handleStyle} />

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

export default PipelineNodeTransformer
