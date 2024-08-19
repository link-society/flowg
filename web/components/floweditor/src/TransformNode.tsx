import React, { useCallback, useContext } from 'react'
import { Handle, Position, Node, NodeProps } from '@xyflow/react'

import { HooksContext } from './context'

export type TransformNode = Node<{
  transformer: string
}>

const TransformNode: React.FC<NodeProps<TransformNode>> = ({ id, data }) => {
  const hooksCtx = useContext(HooksContext)

  const onChange: React.ChangeEventHandler<HTMLInputElement> = useCallback(
    (evt) => {
      hooksCtx.setNodes((nodes) => {
        for (const node of nodes) {
          if (node.id === id) {
            node.data = {transformer: evt.target!.value}
            break
          }
        }

        return [...nodes]
      })
    },
    [id, hooksCtx],
  )

  return (
    <>
      <Handle type="target" position={Position.Left} />
      <div
        className="
          flex flex-row items-center
          z-depth-1 px-3 py-1 gap-2
          light-blue lighten-4 black-text
        "
      >
        <i className="material-icons small">filter_alt</i>
        <div className="input-field">
          <input
            className="nodrag"
            id={`transformer-${id}`}
            type="text"
            defaultValue={data.transformer}
            onChange={onChange}
          />
          <label htmlFor={`transformer-${id}`} className="font-semibold">
            Transformer
          </label>
        </div>
      </div>
      <Handle type="source" position={Position.Right} />
    </>
  )
}

export default TransformNode
