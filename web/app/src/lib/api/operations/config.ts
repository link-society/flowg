import * as request from '@/lib/api/request'

import { StreamConfigModel } from '@/lib/models'

export const listTransformers = async (): Promise<string[]> => {
  type ListTransformersResponse = {
    success: boolean
    transformers: string[]
  }

  const { body } = await request.GET<ListTransformersResponse>('/api/v1/transformers')
  return body.transformers
}

export const listStreams = async (): Promise<{ [stream: string]: StreamConfigModel }> => {
  type ListStreamsResponse = {
    success: boolean
    streams: { [stream: string]: StreamConfigModel }
  }

  const { body } = await request.GET<ListStreamsResponse>('/api/v1/streams')
  return body.streams
}

export const listAlerts = async (): Promise<string[]> => {
  type ListAlertsResponse = {
    success: boolean
    alerts: string[]
  }

  const { body } = await request.GET<ListAlertsResponse>('/api/v1/alerts')
  return body.alerts
}

export const listPipelines = async (): Promise<string[]> => {
  type ListPipelinesResponse = {
    success: boolean
    pipelines: string[]
  }

  const { body } = await request.GET<ListPipelinesResponse>('/api/v1/pipelines')
  return body.pipelines
}
