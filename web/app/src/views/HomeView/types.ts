import StreamConfigModel from '@/lib/models/StreamConfigModel'

export type HomeViewData = {
  streams: { [stream: string]: StreamConfigModel } | null
  transformers: string[] | null
  forwarders: string[] | null
  pipelines: string[] | null
}
