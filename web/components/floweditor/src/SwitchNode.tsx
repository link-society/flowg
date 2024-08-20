import React, { useCallback, useContext } from 'react'
import { Handle, Position, Node, NodeProps } from '@xyflow/react'

import { HooksContext } from './context'

export type SwitchNode = Node<{
  condition: string
}>

const SwitchNode: React.FC<NodeProps<SwitchNode>> = ({ id, data }) => {
  const hooksCtx = useContext(HooksContext)

  const onChange: React.ChangeEventHandler<HTMLInputElement> = useCallback(
    (evt) => {
      hooksCtx.setNodes((nodes) => {
        for (const node of nodes) {
          if (node.id === id) {
            node.data = {condition: evt.target!.value}
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
          red lighten-4 black-text
          hoverable
        "
      >
        <i className="material-icons small">device_hub</i>
        <div className="input-field">
          <input
            className="nodrag"
            id={`switch-${id}`}
            type="text"
            defaultValue={data.condition}
            onChange={onChange}
          />
          <label htmlFor={`switch-${id}`} className="font-semibold">
            Condition
          </label>
        </div>
      </div>
      <Handle type="source" position={Position.Right} />
    </>
  )
}

export default SwitchNode
