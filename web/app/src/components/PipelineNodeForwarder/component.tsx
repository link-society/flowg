import { useTranslation } from 'react-i18next'

import TextField from '@mui/material/TextField'

import ForwardToInboxIcon from '@mui/icons-material/ForwardToInbox'

import { Handle, NodeProps, Position } from '@xyflow/react'

import DialogForwarderEditor from '@/components/DialogForwarderEditor/component'
import PipelineDeleteNodeButton from '@/components/PipelineDeleteNodeButton/component'
import PipelineTraceNodeButton from '@/components/PipelineTraceNodeButton/component'
import PipelineTraceNodeIndicator from '@/components/PipelineTraceNodeIndicator/component'

import { NodeBody, NodeIcon, NodeRoot, ToolbarRow, handleStyle } from './styles'
import { PipelineNodeForwarderData } from './types'

const PipelineNodeForwarder = ({
  id,
  data,
  selected,
}: NodeProps<PipelineNodeForwarderData>) => {
  const { t } = useTranslation()

  return (
    <>
      {selected && (
        <ToolbarRow>
          <DialogForwarderEditor forwarderName={data.forwarder} />
          <PipelineDeleteNodeButton nodeId={id} />
          {data.traces && <PipelineTraceNodeButton traces={data.traces} />}
        </ToolbarRow>
      )}

      <Handle type="target" position={Position.Left} style={handleStyle} />
      <NodeRoot>
        <NodeIcon>
          <ForwardToInboxIcon />
        </NodeIcon>
        <NodeBody className="nodrag">
          <TextField
            label={t('components.pipelineNodeForwarder.label')}
            type="text"
            value={data.forwarder}
            slotProps={{
              input: {
                readOnly: true,
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

export default PipelineNodeForwarder
