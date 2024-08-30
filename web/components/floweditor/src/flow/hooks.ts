import React, { createContext } from 'react'
import { Node } from '@xyflow/react'

export default createContext<{
  setNodes: React.Dispatch<React.SetStateAction<Node[]>>
}>({
  setNodes: () => {},
})
