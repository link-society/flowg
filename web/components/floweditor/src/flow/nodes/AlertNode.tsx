import React, { useCallback, useContext, useEffect } from 'react'
import { Handle, Position, Node, NodeProps, NodeToolbar } from '@xyflow/react'

import HooksContext from '../hooks'

export type AlertNode = Node<{
  alert: string
}>

const AlertNode: React.FC<NodeProps<AlertNode>> = ({ id, data, selected }) => {
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
            node.data = {alert: evt.target!.value}
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
      {data.alert
        ? (
          <NodeToolbar>
            <a
              href={`/web/alerts/edit/${data.alert}/`}
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
          border: '4px solid #2E7D32',
        }}
      >
        <div className="green white-text px-3 py-1 flex flex-row items-center">
          <i className="material-icons small">notifications_active</i>
        </div>
        <div className="input-field px-3 py-1">
          <input
            className="nodrag"
            id={`alert-${id}`}
            type="text"
            defaultValue={data.alert}
            onChange={onChange}
          />
          <label htmlFor={`alert-${id}`} className="font-semibold">
            Alert
          </label>
        </div>
      </div>
    </>
  )
}

export default AlertNode
