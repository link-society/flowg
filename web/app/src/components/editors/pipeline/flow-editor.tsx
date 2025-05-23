import { v4 as uuidv4 } from 'uuid'

import React, {
  DragEventHandler,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from 'react'

import Chip from '@mui/material/Chip'
import Paper from '@mui/material/Paper'
import * as colors from '@mui/material/colors'

import DeviceHubIcon from '@mui/icons-material/DeviceHub'

import {
  Background,
  Controls,
  type Edge,
  type Node,
  type OnConnect,
  type OnEdgesChange,
  type OnNodesChange,
  Panel,
  ReactFlow,
  addEdge,
  applyEdgeChanges,
  applyNodeChanges,
  useReactFlow,
} from '@xyflow/react'

import { PipelineModel } from '@/lib/models/pipeline'

import { HooksContext } from './hooks'
import { ForwarderNode } from './nodes/forwarder'
import { PipelineNode } from './nodes/pipeline'
import { RouterNode } from './nodes/router'
import { SourceNode } from './nodes/source'
import { SwitchNode } from './nodes/switch'
import { TransformNode } from './nodes/transform'

type FlowEditorProps = Readonly<{
  flow: PipelineModel
  onFlowChange: (flow: PipelineModel) => void
}>

export const FlowEditor: React.FC<FlowEditorProps> = ({
  flow,
  onFlowChange,
}) => {
  const { screenToFlowPosition } = useReactFlow()

  const nodeTypes = useMemo(
    () => ({
      source: SourceNode,
      transform: TransformNode,
      switch: SwitchNode,
      pipeline: PipelineNode,
      forwarder: ForwarderNode,
      router: RouterNode,
    }),
    []
  )

  const [nodes, setNodes] = useState<Node[]>(
    flow.nodes.map((node) => {
      if (node.type === 'source') {
        node.deletable = false
      }

      return node
    })
  )

  const [edges, setEdges] = useState<Edge[]>(flow.edges)

  const hooksCtx = useContext(HooksContext)
  hooksCtx.setNodes = setNodes

  useEffect(() => {
    const initialNodes = flow.nodes
    const initialEdges = flow.edges

    setNodes((oldNodes) => {
      const oldNodesById = oldNodes.reduce(
        (acc, node) => {
          acc[node.id] = node
          return acc
        },
        {} as { [key: string]: Node }
      )

      return initialNodes.map((node: Node) => {
        const oldNode = oldNodesById[node.id]
        if (oldNode !== undefined && oldNode.measured) {
          node.measured = oldNode.measured
        }

        return node
      })
    })
    setEdges(initialEdges)
  }, [flow])

  useEffect(() => {
    onFlowChange({ nodes, edges })
  }, [nodes, edges])

  const onNodesChange: OnNodesChange = (changes) =>
    setNodes((nds) => applyNodeChanges(changes, nds))

  const onEdgesChange: OnEdgesChange = (changes) =>
    setEdges((eds) => applyEdgeChanges(changes, eds))

  const onConnect: OnConnect = (conn) => setEdges((eds) => addEdge(conn, eds))

  const onDragOver: DragEventHandler<HTMLDivElement> = (event) => {
    event.preventDefault()
    event.dataTransfer.dropEffect = 'move'
  }

  const onDrop: DragEventHandler<HTMLDivElement> = useCallback(
    (event) => {
      event.preventDefault()

      const itemToNodeMap = {
        transformer: 'transform',
        stream: 'router',
        forwarder: 'forwarder',
        pipeline: 'pipeline',
        switch: 'switch',
      }
      const itemType = event.dataTransfer.getData(
        'item-type'
      ) as keyof typeof itemToNodeMap
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
          id: `node-${uuidv4()}`,
          type,
          position,
          data: {},
        }

        switch (type) {
          case 'transform':
            newNode.data = { transformer: event.dataTransfer.getData('item') }
            break

          case 'router':
            newNode.data = { stream: event.dataTransfer.getData('item') }
            break

          case 'forwarder':
            newNode.data = { forwarder: event.dataTransfer.getData('item') }
            break

          case 'pipeline':
            newNode.data = { pipeline: event.dataTransfer.getData('item') }
            break
        }

        return [...nds, newNode]
      })
    },
    [screenToFlowPosition]
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
          <Paper
            variant="outlined"
            className="flex flex-row items-center gap-3 shadow-xs"
          >
            <div
              className="
                self-stretch flex flex-row items-center
                text-xs font-semibold
                bg-gray-100 border-r border-r-gray-200
                p-2
              "
            >
              <span>Other Nodes:</span>
            </div>

            <div
              className="
                flex flex-row items-center gap-2
                p-2
              "
            >
              <Chip
                icon={<DeviceHubIcon />}
                label="switch"
                variant="outlined"
                sx={{
                  backgroundColor: colors.red[50],
                  borderColor: colors.red[500],
                }}
                className="rounded-none! shadow-xs hover:shadow-lg font-mono!"
                draggable
                onDragStart={(evt) => {
                  evt.dataTransfer.setData('item-type', 'switch')
                  evt.dataTransfer.effectAllowed = 'move'
                }}
              />
            </div>
          </Paper>
        </Panel>
      </ReactFlow>
    </Paper>
  )
}

export default FlowEditor
