import React, { useState } from 'react'
import { ReactFlow } from '@xyflow/react'

import '@xyflow/react/dist/style.css'

const FlowEditor: React.FC = () => {
  const [nodes, setNodes] = useState([
    { id: '1', position: { x: 0, y: 0 }, data: { label: '1' } },
    { id: '2', position: { x: 0, y: 100 }, data: { label: '2' } },
  ])
  const [edges, setEdges] = useState([
    { id: 'e1-2', source: '1', target: '2' },
  ])

  return (
    <div className="w-full h-full">
      <ReactFlow
        nodes={nodes}
        edges={edges}
      />
    </div>
  )
}

export default FlowEditor
