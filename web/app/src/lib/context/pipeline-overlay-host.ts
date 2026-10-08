import { createContext } from 'react'

const PipelineOverlayHostContext = createContext<HTMLDivElement | null>(null)

export default PipelineOverlayHostContext
