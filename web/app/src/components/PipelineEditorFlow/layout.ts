import dagre from 'dagre'

import type { Edge, Node } from '@xyflow/react'

const DEFAULT_NODE_WIDTH = 270
const DEFAULT_NODE_HEIGHT = 100

export type LayoutDirection = 'TB' | 'LR'

export function getLayoutedNodes(
  nodes: Node[],
  edges: Edge[],
  direction: LayoutDirection = 'LR'
): Node[] {
  const graph = new dagre.graphlib.Graph()
  graph.setDefaultEdgeLabel(() => ({}))
  graph.setGraph({ rankdir: direction, nodesep: 60, ranksep: 80 })

  for (const node of nodes) {
    graph.setNode(node.id, {
      width: node.measured?.width ?? DEFAULT_NODE_WIDTH,
      height: node.measured?.height ?? DEFAULT_NODE_HEIGHT,
    })
  }

  for (const edge of edges) {
    graph.setEdge(edge.source, edge.target)
  }

  dagre.layout(graph)

  return nodes.map((node) => {
    const { x: left, y: top } = graph.node(node.id)
    const width = node.measured?.width ?? DEFAULT_NODE_WIDTH
    const height = node.measured?.height ?? DEFAULT_NODE_HEIGHT

    const centerX = left + width / 2
    const centerY = top + height / 2

    return {
      ...node,
      position: {
        x: centerX,
        y: centerY,
      },
    }
  })
}
