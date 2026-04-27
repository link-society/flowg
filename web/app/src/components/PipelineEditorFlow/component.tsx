import { v4 as uuidv4 } from 'uuid'

import React, {
  DragEventHandler,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from 'react'

import Chip from '@mui/material/Chip'
import Typography from '@mui/material/Typography'
import * as colors from '@mui/material/colors'

import DeviceHubIcon from '@mui/icons-material/DeviceHub'

import {
  Background,
  Controls,
  type Edge,
  type KeyCode,
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

import PipelineModel from '@/lib/models/PipelineModel'
import { PipelineTrace } from '@/lib/models/PipelineTrace.ts'

import PipelineEditorHooksProvider from '@/components/PipelineEditorHooksProvider/component'
import PipelineNodeForwarder from '@/components/PipelineNodeForwarder/component'
import PipelineNodePipeline from '@/components/PipelineNodePipeline/component'
import PipelineNodeRouter from '@/components/PipelineNodeRouter/component'
import PipelineNodeSource from '@/components/PipelineNodeSource/component'
import PipelineNodeSwitch from '@/components/PipelineNodeSwitch/component'
import PipelineNodeTransformer from '@/components/PipelineNodeTransformer/component'

import {
  FlowPanelChips,
  FlowPanelLabel,
  FlowPanelPaper,
  FlowRoot,
} from './styles'

type ShortcutMap = {
  deleteKeyCode: KeyCode
  selectionKeyCode: KeyCode
  panActivationKeyCode: KeyCode
  multiSelectionKeyCode: KeyCode
  zoomActivationKeyCode: KeyCode
}

const shortcuts: ShortcutMap = {
  deleteKeyCode: ['Backspace', 'Delete'],
  selectionKeyCode: 'Shift',
  panActivationKeyCode: 'Space',
  multiSelectionKeyCode: ['Meta', 'Control'],
  zoomActivationKeyCode: ['Meta', 'Control'],
}

type PipelineEditorFlowProps = Readonly<{
  flow: PipelineModel
  onFlowChange: (flow: PipelineModel) => void
  pipelineTrace: PipelineTrace | null
}>

export const PipelineEditorFlow: React.FC<PipelineEditorFlowProps> = ({
  flow,
  onFlowChange,
  pipelineTrace,
}) => {
  const { screenToFlowPosition } = useReactFlow()

  const nodeTypes = useMemo(
    () => ({
      source: PipelineNodeSource,
      transform: PipelineNodeTransformer,
      switch: PipelineNodeSwitch,
      pipeline: PipelineNodePipeline,
      forwarder: PipelineNodeForwarder,
      router: PipelineNodeRouter,
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
    setNodes(
      nodes.map((node) => ({
        ...node,
        data: {
          ...node.data,
          trace: pipelineTrace?.find((trace) => trace.nodeID == node.id),
        },
      }))
    )
  }, [pipelineTrace])

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
    <FlowRoot>
      <PipelineEditorHooksProvider value={{ setNodes }}>
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
          {...shortcuts}
        >
          <Background />
          <Controls />

          <Panel position="top-left">
            <FlowPanelPaper variant="outlined">
              <FlowPanelLabel>
                <Typography variant="text">Other Nodes:</Typography>
              </FlowPanelLabel>

              <FlowPanelChips>
                <Chip
                  icon={<DeviceHubIcon />}
                  label="switch"
                  variant="outlined"
                  sx={{
                    backgroundColor: colors.red[50],
                    borderColor: colors.red[500],
                    borderRadius: 0,
                    boxShadow: 1,
                    fontFamily: 'monospace',
                    '&:hover': { boxShadow: 4 },
                  }}
                  draggable
                  onDragStart={(evt) => {
                    evt.dataTransfer.setData('item-type', 'switch')
                    evt.dataTransfer.effectAllowed = 'move'
                  }}
                />
              </FlowPanelChips>
            </FlowPanelPaper>
          </Panel>
        </ReactFlow>
      </PipelineEditorHooksProvider>
    </FlowRoot>
  )
}

export default PipelineEditorFlow
