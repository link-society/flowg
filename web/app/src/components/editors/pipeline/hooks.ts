import React, { createContext } from 'react'
import { Node } from '@xyflow/react'

export const HooksContext = createContext<{
  setNodes: React.Dispatch<React.SetStateAction<Node[]>>
}>({
  setNodes: () => {},
})
