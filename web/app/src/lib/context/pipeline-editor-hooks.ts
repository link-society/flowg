import React, { createContext } from 'react'

import { Node } from '@xyflow/react'

const PipelineEditorHooksContext = createContext<{
  setNodes: React.Dispatch<React.SetStateAction<Node[]>>
}>({
  setNodes: () => {},
})

export default PipelineEditorHooksContext
