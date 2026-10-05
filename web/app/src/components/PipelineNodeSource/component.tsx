import { useTranslation } from 'react-i18next'

import TextField from '@mui/material/TextField'
import { useTheme } from '@mui/material/styles'

import InputIcon from '@mui/icons-material/Input'

import { Handle, NodeProps, Position } from '@xyflow/react'

import {
  NodeBody,
  NodeIcon,
  NodeRoot,
  ToolbarRow,
  handleStyle,
} from '@/components/PipelineNodeCard/styles'
import PipelineTraceNodeButton from '@/components/PipelineTraceNodeButton/component'
import PipelineTraceNodeIndicator from '@/components/PipelineTraceNodeIndicator/component'

import { PipelineNodeSourceData } from './types'

const PipelineNodeSource = ({
  selected,
  data,
}: NodeProps<PipelineNodeSourceData>) => {
  const { t } = useTranslation()
  const theme = useTheme()

  return (
    <>
      {selected && data.traces && (
        <ToolbarRow>
          <PipelineTraceNodeButton traces={data.traces} />
        </ToolbarRow>
      )}

      <NodeRoot borderColor={theme.tokens.colors.nodeSourceBorder}>
        <NodeIcon bgColor={theme.tokens.colors.nodeSourceBg}>
          <InputIcon />
        </NodeIcon>
        <NodeBody className="nodrag">
          <TextField
            label={t('components.pipelineNodeSource.label')}
            type="text"
            value={data.type.toUpperCase()}
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

export default PipelineNodeSource
