import { useContext } from 'react'

import PipelineEditorHooksContext from '@/lib/context/pipeline-editor-hooks'

export const usePipelineEditorHooks = () => {
  return useContext(PipelineEditorHooksContext)
}
