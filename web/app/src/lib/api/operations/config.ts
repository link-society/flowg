import * as request from '@/lib/api/request'

import { PipelineModel, StreamConfigModel } from '@/lib/models'

export const listTransformers = async (): Promise<string[]> => {
  type ListTransformersResponse = {
    success: boolean
    transformers: string[]
  }

  const { body } = await request.GET<ListTransformersResponse>({
    path: '/api/v1/transformers',
  })
  return body.transformers
}

export const getTransformer = async (transformer: string): Promise<string> => {
  type GetTransformerResponse = {
    success: boolean
    script: string
  }

  const { body } = await request.GET<GetTransformerResponse>({
    path: `/api/v1/transformers/${transformer}`,
  })

  return body.script
}

export const saveTransformer = async (transformer: string, script: string): Promise<void> => {
  type SaveTransformerRequest = {
    script: string
  }

  type SaveTransformerResponse = {
    success: boolean
  }

  await request.PUT<SaveTransformerRequest, SaveTransformerResponse>({
    path: `/api/v1/transformers/${transformer}`,
    body: { script },
  })
}

export const deleteTransformer = async (transformer: string): Promise<void> => {
  type DeleteTransformerResponse = {
    success: boolean
  }

  await request.DELETE<DeleteTransformerResponse>({
    path: `/api/v1/transformers/${transformer}`,
  })
}

export const listStreams = async (): Promise<{ [stream: string]: StreamConfigModel }> => {
  type ListStreamsResponse = {
    success: boolean
    streams: { [stream: string]: StreamConfigModel }
  }

  const { body } = await request.GET<ListStreamsResponse>({
    path: '/api/v1/streams',
  })
  return body.streams
}

export const listStreamFields = async (stream: string): Promise<string[]> => {
  type ListStreamFieldsResponse = {
    success: boolean
    fields: string[]
  }

  const { body } = await request.GET<ListStreamFieldsResponse>({
    path: `/api/v1/streams/${stream}/fields`,
  })

  return body.fields
}

export const configureStream = async (stream: string, config: StreamConfigModel): Promise<void> => {
  type ConfigureStreamRequest = {
    config: StreamConfigModel
  }

  type ConfigureStreamResponse = {
    success: boolean
  }

  await request.PUT<ConfigureStreamRequest, ConfigureStreamResponse>({
    path: `/api/v1/streams/${stream}`,
    body: { config },
  })
}

export const purgeStream = async (stream: string): Promise<void> => {
  type PurgeStreamResponse = {
    success: boolean
  }

  await request.DELETE<PurgeStreamResponse>({
    path: `/api/v1/streams/${stream}`,
  })
}

export const listAlerts = async (): Promise<string[]> => {
  type ListAlertsResponse = {
    success: boolean
    alerts: string[]
  }

  const { body } = await request.GET<ListAlertsResponse>({
    path: '/api/v1/alerts',
  })
  return body.alerts
}

export const listPipelines = async (): Promise<string[]> => {
  type ListPipelinesResponse = {
    success: boolean
    pipelines: string[]
  }

  const { body } = await request.GET<ListPipelinesResponse>({
    path: '/api/v1/pipelines',
  })

  return body.pipelines
}

export const getPipeline = async (pipeline: string): Promise<PipelineModel> => {
  type GetPipelineResponse = {
    success: boolean
    flow: PipelineModel
  }

  const { body } = await request.GET<GetPipelineResponse>({
    path: `/api/v1/pipelines/${pipeline}`,
  })
  return body.flow
}

export const savePipeline = async (pipeline: string, flow: PipelineModel): Promise<void> => {
  type SavePipelineRequest = {
    flow: PipelineModel
  }

  type SavePipelineResponse = {
    success: boolean
  }

  await request.PUT<SavePipelineRequest, SavePipelineResponse>({
    path: `/api/v1/pipelines/${pipeline}`,
    body: { flow },
  })
}

export const deletePipeline = async (pipeline: string): Promise<void> => {
  type DeletePipelineResponse = {
    success: boolean
  }

  await request.DELETE<DeletePipelineResponse>({
    path: `/api/v1/pipelines/${pipeline}`,
  })
}
