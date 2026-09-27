import type { FlowKind } from './config'

export type LogFlowDiagramProps = Readonly<{ labelledBy: string }>
export type EventFlowProps = Readonly<{ compact?: boolean }>
export type EventShapeProps = Readonly<{ kind: FlowKind['name'] }>
