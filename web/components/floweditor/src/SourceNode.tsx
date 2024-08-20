import React from 'react'
import { Handle, Position, NodeProps } from '@xyflow/react'


const SourceNode: React.FC<NodeProps> = ({}) => {
  return (
    <>
      <div
        className="
          flex flex-row items-center
          z-depth-1 px-3 py-1 gap-2
          orange lighten-1 black-text
          hoverable
        "
      >
        <i className="material-icons small">input</i>
        <span className="font-semibold">Log Source</span>
      </div>
      <Handle type="source" position={Position.Right} />
    </>
  )
}

export default SourceNode
