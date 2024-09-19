import React, {
  DragEventHandler,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from 'react'

import * as colors from '@mui/material/colors'

import DeviceHubIcon from '@mui/icons-material/DeviceHub'

import Paper from '@mui/material/Paper'
import Chip from '@mui/material/Chip'

import {
  ReactFlow,
  Background,
  Controls,
  Panel,
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

import { HooksContext } from './hooks'

import { SourceNode } from './nodes/source'
import { TransformNode } from './nodes/transform'
import { SwitchNode } from './nodes/switch'
import { PipelineNode } from './nodes/pipeline'
import { AlertNode } from './nodes/alert'
import { RouterNode } from './nodes/router'

import { PipelineModel } from '@/lib/models'

type FlowEditorProps = {
  flow: PipelineModel
  onFlowChange: (flow: PipelineModel) => void
}

export const FlowEditor: React.FC<FlowEditorProps> = ({ flow, onFlowChange }) => {
  const { screenToFlowPosition } = useReactFlow()

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

  const [nodes, setNodes] = useState<Node[]>(flow.nodes.map(node => {
    if (node.type === 'source') {
      node.deletable = false
    }

    return node
  }))

  const [edges, setEdges] = useState<Edge[]>(flow.edges)

  const hooksCtx = useContext(HooksContext)
  hooksCtx.setNodes = setNodes

  useEffect(
    () => {
      const initialNodes = flow.nodes
      const initialEdges = flow.edges

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
    () => { onFlowChange({ nodes, edges }) },
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

      const itemToNodeMap = {
        transformer: 'transform',
        stream: 'router',
        alert: 'alert',
        pipeline: 'pipeline',
        switch: 'switch',
      }
      const itemType = event.dataTransfer.getData('item-type') as keyof typeof itemToNodeMap
      const dndNodeType = itemToNodeMap[itemType]

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
            newNode.data = {transformer: event.dataTransfer.getData('item')}
            break

          case 'router':
            newNode.data = {stream: event.dataTransfer.getData('item')}
            break

          case 'alert':
            newNode.data = {alert: event.dataTransfer.getData('item')}
            break

          case 'pipeline':
            newNode.data = {pipeline: event.dataTransfer.getData('item')}
            break
        }

        return [...nds, newNode]
      })
    },
    [screenToFlowPosition, setNodes],
  )

  return (
    <Paper className="w-full h-full">
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

        <Panel position="top-left">
          <Paper variant="outlined" className="p-2 flex flex-row items-center gap-2">
            <span className="text-xs font-semibold">Special Nodes:</span>

            <Chip
              icon={<DeviceHubIcon />}
              label="switch"
              variant="outlined"
              sx={{
                backgroundColor: colors.red[50],
                borderColor: colors.red[500],
              }}
              className="!rounded-none shadow-sm hover:shadow-lg !font-mono"
              draggable
              onDragStart={(evt) => {
                evt.dataTransfer.setData('item-type', 'switch')
                evt.dataTransfer.effectAllowed = 'move'
              }}
            />
          </Paper>
        </Panel>
      </ReactFlow>
    </Paper>
  )
}

export default FlowEditor
