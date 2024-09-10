import React, { DragEventHandler, useCallback, useContext, useEffect, useMemo, useState } from 'react'
import {
  ReactFlow,
  Background,
  Controls,
  addEdge,
  applyEdgeChanges,
  applyNodeChanges,
  useReactFlow,
  type Node,
  type Edge,
  type OnNodesChange,
  type OnEdgesChange,
  type OnConnect,
} from '@xyflow/react'

import '@xyflow/react/dist/style.css'

import HooksContext from './hooks'

import NodeSelector from './NodeSelector'

import SourceNode from './nodes/SourceNode'
import TransformNode from './nodes/TransformNode'
import SwitchNode from './nodes/SwitchNode'
import PipelineNode from './nodes/PipelineNode'
import AlertNode from './nodes/AlertNode'
import RouterNode from './nodes/RouterNode'
import { useDnD } from '../dnd/context'

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
}

const FlowEditor: React.FC<FlowEditorProps> = ({ flow, onFlowChange }) => {
  const { screenToFlowPosition } = useReactFlow()
  const [dndNodeType] = useDnD()

  const nodeTypes = useMemo(
    () => ({
      source: SourceNode,
      transform: TransformNode,
      switch: SwitchNode,
      pipeline: PipelineNode,
      alert: AlertNode,
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

      setNodes((oldNodes) => {
        const oldNodesById = oldNodes.reduce((acc, node) => {
          acc[node.id] = node
          return acc
        }, {} as { [key: string]: Node })

        return initialNodes.map((node: Node) => {
          const oldNode = oldNodesById[node.id]
          if (oldNode !== undefined && oldNode.measured) {
            node.measured = oldNode.measured
          }

          return node
        })
      })
      setEdges(initialEdges)
    },
    [flow],
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

  const onDragOver: DragEventHandler<HTMLDivElement> = useCallback(
    (event) => {
      event.preventDefault()
      event.dataTransfer.dropEffect = 'move'
    },
    [],
  )

  const onDrop: DragEventHandler<HTMLDivElement> = useCallback(
    (event) => {
      event.preventDefault()

      if (!dndNodeType) {
        return
      }

      const type = dndNodeType
      const position = screenToFlowPosition({
        x: event.clientX,
        y: event.clientY,
      })

      setNodes((nds) => {
        const newNode = {
          id: `node-${nds.length}`,
          type,
          position,
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
    },
    [screenToFlowPosition, setNodes, dndNodeType],
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
        onDrop={onDrop}
        onDragOver={onDragOver}
        fitView
        snapToGrid
        defaultEdgeOptions={{ animated: true, type: 'smoothstep' }}
      >
        <Background />
        <Controls />
        <NodeSelector />
      </ReactFlow>
    </div>
  )
}

export default FlowEditor
