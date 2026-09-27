import { useTranslation } from 'react-i18next'

import DragIndicatorIcon from '@mui/icons-material/DragIndicator'

import { DragHandleRoot } from './styles'

const PIPELINE_NODE_DRAG_HANDLE_CLASS = 'pipeline-node-drag-handle'
export const PIPELINE_NODE_DRAG_HANDLE_SELECTOR = `.${PIPELINE_NODE_DRAG_HANDLE_CLASS}`

const PipelineNodeDragHandle = () => {
  const { t } = useTranslation()

  return (
    <DragHandleRoot
      className={PIPELINE_NODE_DRAG_HANDLE_CLASS}
      onClick={(evt) => evt.stopPropagation()}
      aria-label={t('components.pipelineNodeDragHandle.label')}
    >
      <DragIndicatorIcon fontSize="small" />
    </DragHandleRoot>
  )
}

export default PipelineNodeDragHandle
