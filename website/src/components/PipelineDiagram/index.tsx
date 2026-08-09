import React, { useEffect, useRef, useState } from 'react'

import styles from './styles.module.css'

/** Slightly wider than the editor's 270px, so field values fit once scaled down. */
export const NODE_WIDTH = 300
export const NODE_HEIGHT = 100

export type PipelineNodeType =
  | 'source'
  | 'transformer'
  | 'switch'
  | 'metric'
  | 'router'
  | 'forwarder'
  | 'pipeline'

type NodeKind = {
  bg: string
  border: string
  label: string
  icon: string
  hasTarget: boolean
  hasSource: boolean
}

/** Mirrors web/app/src/theme/tokens/colors.ts and the icons of each node type. */
const KINDS: Record<PipelineNodeType, NodeKind> = {
  source: {
    bg: '#f97316',
    border: '#c2410c',
    label: 'Source',
    icon: 'M21 3.01H3c-1.1 0-2 .9-2 2V9h2V4.99h18v14.03H3V15H1v4.01c0 1.1.9 1.98 2 1.98h18c1.1 0 2-.88 2-1.98v-14c0-1.11-.9-2-2-2M11 16l4-4-4-4v3H1v2h10z',
    hasTarget: false,
    hasSource: true,
  },
  transformer: {
    bg: '#1d4ed8',
    border: '#1e3a8a',
    label: 'Transformer',
    icon: 'M4.25 5.61C6.27 8.2 10 13 10 13v6c0 .55.45 1 1 1h2c.55 0 1-.45 1-1v-6s3.72-4.8 5.74-7.39c.51-.66.04-1.61-.79-1.61H5.04c-.83 0-1.3.95-.79 1.61',
    hasTarget: true,
    hasSource: true,
  },
  switch: {
    bg: '#dc2626',
    border: '#b91c1c',
    label: 'Condition',
    icon: 'm17 16-4-4V8.82C14.16 8.4 15 7.3 15 6c0-1.66-1.34-3-3-3S9 4.34 9 6c0 1.3.84 2.4 2 2.82V12l-4 4H3v5h5v-3.05l4-4.2 4 4.2V21h5v-5z',
    hasTarget: true,
    hasSource: true,
  },
  metric: {
    bg: '#0a3c36',
    border: '#004d40',
    label: 'Name',
    icon: 'M4 9h4v11H4zm12 4h4v7h-4zm-6-9h4v16h-4z',
    hasTarget: true,
    hasSource: false,
  },
  router: {
    bg: '#7e22ce',
    border: '#581c87',
    label: 'Stream',
    icon: 'M2 20h20v-4H2zm2-3h2v2H4zM2 4v4h20V4zm4 3H4V5h2zm-4 7h20v-4H2zm2-3h2v2H4z',
    hasTarget: true,
    hasSource: false,
  },
  forwarder: {
    bg: '#15803d',
    border: '#14532d',
    label: 'Forwarder',
    icon: 'M20 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h9v-2H4V8l8 5 8-5v5h2V6c0-1.1-.9-2-2-2m-8 7L4 6h16zm7 4 4 4-4 4v-3h-4v-2h4z',
    hasTarget: true,
    hasSource: false,
  },
  pipeline: {
    bg: '#eab308',
    border: '#ca8a04',
    label: 'Pipeline',
    icon: 'M22 11V3h-7v3H9V3H2v8h7V8h2v10h4v3h7v-8h-7v3h-2V8h2v3z',
    hasTarget: true,
    hasSource: false,
  },
}

export type PipelineNodeSpec = {
  id: string
  type: PipelineNodeType
  x: number
  y: number
  value: string
  label?: string
}

export type PipelineEdgeSpec = {
  from: string
  to: string
}

type PipelineNodeProps = {
  type: PipelineNodeType
  value: string
  label?: string
  style?: React.CSSProperties
}

export function PipelineNode({ type, value, label, style }: PipelineNodeProps) {
  const kind = KINDS[type]

  return (
    <div
      className={styles.node}
      style={
        {
          '--node-bg': kind.bg,
          '--node-border': kind.border,
          width: NODE_WIDTH,
          height: NODE_HEIGHT,
          ...style,
        } as React.CSSProperties
      }
    >
      {kind.hasTarget && (
        <span className={`${styles.handle} ${styles.handleTarget}`} />
      )}

      <div className={styles.icon}>
        <svg viewBox="0 0 24 24" width="24" height="24" fill="currentColor">
          <path d={kind.icon} />
        </svg>
      </div>

      <div className={styles.body}>
        <div className={styles.field}>
          <span className={styles.fieldLabel}>{label ?? kind.label}</span>
          <span className={styles.fieldValue}>{value}</span>
        </div>
      </div>

      {kind.hasSource && (
        <span className={`${styles.handle} ${styles.handleSource}`} />
      )}
    </div>
  )
}

/** Orthogonal path with rounded corners, like React Flow's `smoothstep` edges. */
function smoothStepPath(
  sx: number,
  sy: number,
  tx: number,
  ty: number,
  radius = 8,
): string {
  if (Math.abs(sy - ty) < 1) {
    return `M ${sx},${sy} L ${tx},${ty}`
  }

  const cx = sx + (tx - sx) / 2
  const dir = ty > sy ? 1 : -1
  const r = Math.min(radius, Math.abs(ty - sy) / 2, Math.abs(cx - sx))

  return [
    `M ${sx},${sy}`,
    `L ${cx - r},${sy}`,
    `Q ${cx},${sy} ${cx},${sy + dir * r}`,
    `L ${cx},${ty - dir * r}`,
    `Q ${cx},${ty} ${cx + r},${ty}`,
    `L ${tx},${ty}`,
  ].join(' ')
}

type PipelineDiagramProps = {
  nodes: PipelineNodeSpec[]
  edges?: PipelineEdgeSpec[]
  padding?: number
  alt?: string
}

export default function PipelineDiagram({
  nodes,
  edges = [],
  padding = 20,
  alt,
}: PipelineDiagramProps) {
  const containerRef = useRef<HTMLDivElement>(null)
  const [scale, setScale] = useState(1)

  const width =
    Math.max(...nodes.map((node) => node.x)) + NODE_WIDTH + padding * 2
  const height =
    Math.max(...nodes.map((node) => node.y)) + NODE_HEIGHT + padding * 2

  useEffect(() => {
    const container = containerRef.current
    if (!container) {
      return
    }

    const observer = new ResizeObserver(([entry]) => {
      setScale(Math.min(1, entry.contentRect.width / width))
    })
    observer.observe(container)
    return () => observer.disconnect()
  }, [width])

  const byId = new Map(nodes.map((node) => [node.id, node]))

  return (
    <div
      ref={containerRef}
      className={styles.container}
      style={{ aspectRatio: `${width} / ${height}` }}
      role="img"
      aria-label={alt}
    >
      <div
        className={styles.canvas}
        style={{ width, height, transform: `scale(${scale})` }}
      >
        <svg className={styles.edges} width={width} height={height}>
          {edges.map(({ from, to }) => {
            const source = byId.get(from)
            const target = byId.get(to)
            if (!source || !target) {
              return null
            }

            return (
              <path
                key={`${from}-${to}`}
                className={styles.edge}
                d={smoothStepPath(
                  source.x + NODE_WIDTH + padding,
                  source.y + NODE_HEIGHT / 2 + padding,
                  target.x + padding,
                  target.y + NODE_HEIGHT / 2 + padding,
                )}
              />
            )
          })}
        </svg>

        {nodes.map((node) => (
          <PipelineNode
            key={node.id}
            type={node.type}
            value={node.value}
            label={node.label}
            style={{ left: node.x + padding, top: node.y + padding }}
          />
        ))}
      </div>
    </div>
  )
}
