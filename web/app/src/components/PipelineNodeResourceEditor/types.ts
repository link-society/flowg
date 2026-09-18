export type PipelineNodeResourceKind = 'stream' | 'transformer' | 'forwarder'

export type PipelineNodeResourceSaveAction = Readonly<{
  saving: boolean
  canSave: boolean
  onSave: () => void
}>

export type PipelineNodeResourceEditorProps = Readonly<{
  kind: PipelineNodeResourceKind
  name: string
  onSaveActionChange: (action: PipelineNodeResourceSaveAction | null) => void
}>
