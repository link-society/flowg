import React from 'react'
import { Handle, Position, Node, NodeProps } from '@xyflow/react'

export type SourceNode = Node<{
  type: string
}>

const SourceNode: React.FC<NodeProps<SourceNode>> = ({ data }) => {
  return (
    <>
      <div
        className="
          flex flex-row items-stretch
          z-depth-1 p-0 gap-2
          white black-text
          hoverable
        "
        style={{
          border: '4px solid #EF6C00',
        }}
      >
        <div className="orange white-text px-3 py-1 flex flex-row items-center">
          <i className="material-icons small">input</i>
        </div>
        <div className="px-3 py-1 flex flex-row items-center">
          <span className="font-semibold">{ data.type.toUpperCase() }</span>
        </div>
      </div>
      <Handle
        type="source"
        position={Position.Right}
        style={{
          width: '12px',
          height: '12px',
        }}
      />
    </>
  )
}

export default SourceNode
