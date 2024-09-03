import React, { useCallback, useContext, useEffect } from 'react'
import { Handle, Position, Node, NodeProps, NodeToolbar } from '@xyflow/react'

import HooksContext from '../hooks'

export type RouterNode = Node<{
  stream: string
}>

const RouterNode: React.FC<NodeProps<RouterNode>> = ({ id, data, selected }) => {
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
      {data.stream
        ? (
          <NodeToolbar>
            <a
              href={`/web/storage/edit/${data.stream}/`}
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
          border: '4px solid #6A1B9A',
        }}
      >
        <div className="purple white-text px-3 py-1 flex flex-row items-center">
          <i className="material-icons small">storage</i>
        </div>
        <div className="input-field px-3 py-1">
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
