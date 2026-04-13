import StreamConfigModel from '@/lib/models/StreamConfigModel'

export type StreamEditorProps = {
  streamConfig: StreamConfigModel
  storageUsage: number
  onStreamConfigChange: (config: StreamConfigModel) => void
}
