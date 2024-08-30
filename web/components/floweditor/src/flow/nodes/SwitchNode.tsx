import React, { useCallback, useContext } from 'react'
import { Handle, Position, Node, NodeProps } from '@xyflow/react'

import HooksContext from '../hooks'

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
      <Handle
        type="target"
        position={Position.Left}
        style={{
          width: '12px',
          height: '12px',
        }}
      />
      <div
        className="
          flex flex-row items-stretch
          z-depth-1 px-0 gap-2
          white black-text
          hoverable
        "
        style={{
          border: '4px solid #C62828',
        }}
      >
        <div className="red darken-2 white-text px-3 py-1 flex flex-row items-center">
          <i className="material-icons small">device_hub</i>
        </div>
        <div className="input-field px-3 py-1">
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

export default SwitchNode
