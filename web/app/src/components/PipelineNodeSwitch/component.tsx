import { useTranslation } from 'react-i18next'

import TextField from '@mui/material/TextField'

import DeviceHubIcon from '@mui/icons-material/DeviceHub'

import { Handle, NodeProps, Position } from '@xyflow/react'

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
  const { t } = useTranslation()

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
          <DeviceHubIcon />
        </NodeIcon>
        <NodeBody className="nodrag">
          <TextField
            label={t('components.pipelineNodeSwitch.label')}
            type="text"
            value={data.condition}
            slotProps={{
              input: {
                readOnly: true,
                sx: { fontFamily: 'monospace' },
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

export default PipelineNodeSwitch
