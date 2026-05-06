import StreamConfigModel from '@/lib/models/StreamConfigModel'

export type LoaderData = {
  streams: Record<string, StreamConfigModel>
  usage: number
  currentStream: string
}
