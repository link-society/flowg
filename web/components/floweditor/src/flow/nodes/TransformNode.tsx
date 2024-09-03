import React, { useCallback, useContext, useEffect } from 'react'
import { Handle, Position, Node, NodeProps, NodeToolbar } from '@xyflow/react'

import HooksContext from '../hooks'

export type TransformNode = Node<{
  transformer: string
}>

const TransformNode: React.FC<NodeProps<TransformNode>> = ({ id, data, selected }) => {
  const hooksCtx = useContext(HooksContext)

  useEffect(
    () => {
      const saveBtn = document.getElementById('action_save')

      if (selected) {
        saveBtn?.classList.add('pulse', 'orange')
      }
      else {
        saveBtn?.classList.remove('pulse', 'orange')
      }
    },
    [selected],
  )

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
    {data.transformer
      ? (
        <NodeToolbar>
          <a
            href={`/web/transformers/edit/${data.transformer}/`}
            className="btn-small waves-effect waves-light"
          >
            <i className="material-icons left">build</i>
            Edit
          </a>
        </NodeToolbar>
      )
      : (
        <></>
      )
    }
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
          z-depth-1 p-0 gap-2
          white black-text
          hoverable
        "
        style={{
          border: '4px solid #1565C0',
        }}
      >
        <div className="blue darken-2 white-text px-3 py-1 flex flex-row items-center">
          <i className="material-icons small">filter_alt</i>
        </div>
        <div className="input-field px-3 py-1">
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

export default TransformNode
