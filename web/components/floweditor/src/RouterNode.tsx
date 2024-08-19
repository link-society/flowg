import React, { useCallback, useContext } from 'react'
import { Handle, Position, Node, NodeProps } from '@xyflow/react'

import { HooksContext } from './context'

export type RouterNode = Node<{
  stream: string
}>

const RouterNode: React.FC<NodeProps<RouterNode>> = ({ id, data }) => {
  const hooksCtx = useContext(HooksContext)

  const onChange: React.ChangeEventHandler<HTMLInputElement> = useCallback(
    (evt) => {
      hooksCtx.setNodes((nodes) => {
        for (const node of nodes) {
          if (node.id === id) {
            node.data = {stream: evt.target!.value}
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
          purple lighten-4 black-text
        "
      >
        <i className="material-icons small">storage</i>
        <div className="input-field">
          <input
            className="nodrag"
            id={`router-${id}`}
            type="text"
            defaultValue={data.stream}
            onChange={onChange}
          />
          <label htmlFor={`router-${id}`} className="font-semibold">
            Stream
          </label>
        </div>
      </div>
    </>
  )
}

export default RouterNode
