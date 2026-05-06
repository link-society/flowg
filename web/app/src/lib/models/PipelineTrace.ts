export type PipelineTrace = NodeTrace &
  {
    nodeID: string
  }[]

export type NodeTrace = {
  input?: Record<string, string>
  output?: Array<Record<string, string>>
  error?: string
}
