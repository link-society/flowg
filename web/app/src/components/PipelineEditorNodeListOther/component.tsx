import { DragEvent } from 'react'
import { useTranslation } from 'react-i18next'

import { useTheme } from '@mui/material/styles'

import BarChartIcon from '@mui/icons-material/BarChart'
import DeviceHubIcon from '@mui/icons-material/DeviceHub'

import {
  NodeChip,
  NodeListHeader,
  NodeListItems,
  NodeListRoot,
  NodeListTitle,
} from '@/components/PipelineEditorNodeList/styles'

const onDragStart = (itemType: string) => (evt: DragEvent) => {
  evt.dataTransfer.setData('item-type', itemType)
  evt.dataTransfer.effectAllowed = 'move'
}

const PipelineEditorNodeListOther = () => {
  const { t } = useTranslation()
  const theme = useTheme()

  return (
    <NodeListRoot>
      <NodeListHeader>
        <NodeListTitle>
          {t('components.pipelineEditorNodeListOther.title')}
        </NodeListTitle>
      </NodeListHeader>

      <NodeListItems>
        <NodeChip
          icon={<DeviceHubIcon />}
          label={t('components.pipelineEditorFlow.switchNodeLabel')}
          variant="outlined"
          chipBgColor={theme.tokens.colors.switchChipBg}
          chipBorderColor={theme.tokens.colors.switchChipBorder}
          draggable
          onDragStart={onDragStart('switch')}
        />
        <NodeChip
          icon={<BarChartIcon />}
          label={t('components.pipelineEditorFlow.metricNodeLabel')}
          variant="outlined"
          chipBgColor={theme.tokens.colors.metricChipBg}
          chipBorderColor={theme.tokens.colors.metricChipBorder}
          draggable
          onDragStart={onDragStart('metric')}
        />
      </NodeListItems>
    </NodeListRoot>
  )
}

export default PipelineEditorNodeListOther
