import React, { useCallback, useContext, useEffect, useMemo, useState } from 'react'
import {
  ReactFlow,
  Background,
  Controls,
  addEdge,
  applyEdgeChanges,
  applyNodeChanges,
  type Node,
  type Edge,
  type OnNodesChange,
  type OnEdgesChange,
  type OnConnect,
} from '@xyflow/react'

import '@xyflow/react/dist/style.css'

import { AddNodeEvent } from './event'
import { HooksContext } from './context'

import SourceNode from './SourceNode'
import TransformNode from './TransformNode'
import SwitchNode from './SwitchNode'
import RouterNode from './RouterNode'

const defaultSourceNode: Node = {
  id: '__builtin__source',
  type: 'source',
  position: { x: 0, y: 0 },
  deletable: false,
  data: {},
}

interface FlowEditorProps {
  flow: string
  onFlowChange: (value: string) => void
  eventTarget: EventTarget
}

const FlowEditor: React.FC<FlowEditorProps> = ({ flow, onFlowChange, eventTarget }) => {
  const nodeTypes = useMemo(
    () => ({
      source: SourceNode,
      transform: TransformNode,
      switch: SwitchNode,
      router: RouterNode,
    }),
    [],
  )

  const flowData = JSON.parse(flow) ?? {}
  const initialNodes = flowData.nodes ?? [defaultSourceNode]
  const initialEdges = flowData.edges ?? []

  const [nodes, setNodes] = useState<Node[]>(initialNodes)
  const [edges, setEdges] = useState<Edge[]>(initialEdges)

  const hooksCtx = useContext(HooksContext)
  hooksCtx.setNodes = setNodes

  useEffect(
    () => {
      const flowData = JSON.parse(flow) ?? {}
      const initialNodes = flowData.nodes ?? [defaultSourceNode]
      const initialEdges = flowData.edges ?? []

      setNodes(initialNodes)
      setEdges(initialEdges)
    },
    [flow],
  )

  useEffect(
    () => {
      const handleAddNode = (event: Event) => {
        const type = (event as AddNodeEvent).detail.type

        setNodes((nds) => {
          const newNode = {
            id: `node-${nds.length}`,
            type,
            position: { x: 0, y: 0 },
            data: {},
          }

          switch (type) {
            case 'transform':
              newNode.data = {transformer: ''}
              break

            case 'router':
              newNode.data = {stream: ''}
              break
          }

          return [...nds, newNode]
        })
      }

      eventTarget.addEventListener('add-node', handleAddNode)

      return () => {
        eventTarget.removeEventListener('add-node', handleAddNode)
      }
    },
    [eventTarget, setNodes],
  )

  useEffect(
    () => {
      onFlowChange(JSON.stringify({ nodes, edges }))
    },
    [nodes, edges, onFlowChange],
  )

  const onNodesChange: OnNodesChange = useCallback(
    (changes) => setNodes((nds) => applyNodeChanges(changes, nds)),
    [setNodes],
  )

  const onEdgesChange: OnEdgesChange = useCallback(
    (changes) => setEdges((eds) => applyEdgeChanges(changes, eds)),
    [setEdges],
  )

  const onConnect: OnConnect = useCallback(
    (conn) => setEdges((eds) => addEdge(conn, eds)),
    [setEdges],
  )

  return (
    <div className="w-full h-full">
      <ReactFlow
        nodeTypes={nodeTypes}
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        snapToGrid
        defaultEdgeOptions={{ animated: true, type: 'smoothstep' }}
      >
        <Background />
        <Controls />
      </ReactFlow>
    </div>
  )
}

export default FlowEditor
