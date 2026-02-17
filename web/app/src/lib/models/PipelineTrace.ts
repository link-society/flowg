export type PipelineTrace = NodeTrace &
  {
    nodeID: string
  }[]

export type NodeTrace = {
  input?: Record<string, string>
  output?: Record<string, string>
}
