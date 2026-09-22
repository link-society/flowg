import { useContext } from 'react'

import PipelineOverlayHostContext from '@/lib/context/pipeline-overlay-host'

export const usePipelineOverlayHost = () => {
  return useContext(PipelineOverlayHostContext)
}
