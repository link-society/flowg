import React from 'react'
import { Panel } from '@xyflow/react'
import { useDnD } from '../dnd/context'

const NodeSelector: React.FC = () => {
  const [_, setType] = useDnD()

  const onDragStart = (event: React.DragEvent<HTMLButtonElement>, nodeType: string) => {
    setType(nodeType)
    event.dataTransfer!.effectAllowed = 'move'
  }

  return (
    <Panel position="bottom-center">
      <div className="flex flex-row items-center gap-2 white z-depth-2 p-1">
        <button
          className="btn-small tooltipped blue"
          data-position="top"
          data-tooltip="Transform Node"
          draggable
          onDragStart={(event) => onDragStart(event, 'transform')}
        >
          <i className="material-icons">filter_alt</i>
        </button>
        <button
          className="btn-small tooltipped red"
          data-position="top"
          data-tooltip="Switch Node"
          draggable
          onDragStart={(event) => onDragStart(event, 'switch')}
        >
          <i className="material-icons">device_hub</i>
        </button>
        <button
          className="btn-small tooltipped purple"
          data-position="top"
          data-tooltip="Router Node"
          draggable
          onDragStart={(event) => onDragStart(event, 'router')}
        >
          <i className="material-icons">storage</i>
        </button>
      </div>
    </Panel>
  )
}

export default NodeSelector
