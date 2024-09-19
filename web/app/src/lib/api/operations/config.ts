import * as request from '@/lib/api/request'

import { PipelineModel, StreamConfigModel } from '@/lib/models'

export const listTransformers = async (): Promise<string[]> => {
  type ListTransformersResponse = {
    success: boolean
    transformers: string[]
  }

  const { body } = await request.GET<ListTransformersResponse>('/api/v1/transformers')
  return body.transformers
}

export const getTransformer = async (transformer: string): Promise<string> => {
  type GetTransformerResponse = {
    success: boolean
    script: string
  }

  const { body } = await request.GET<GetTransformerResponse>(
    `/api/v1/transformers/${transformer}`
  )
  return body.script
}

export const saveTransformer = async (transformer: string, script: string): Promise<void> => {
  type SaveTransformerRequest = {
    script: string
  }

  type SaveTransformerResponse = {
    success: boolean
  }

  await request.PUT<SaveTransformerRequest, SaveTransformerResponse>(
    `/api/v1/transformers/${transformer}`,
    { script }
  )
}

export const deleteTransformer = async (transformer: string): Promise<void> => {
  type DeleteTransformerResponse = {
    success: boolean
  }

  await request.DELETE<DeleteTransformerResponse>(`/api/v1/transformers/${transformer}`)
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

export const getPipeline = async (pipeline: string): Promise<PipelineModel> => {
  type GetPipelineResponse = {
    success: boolean
    flow: PipelineModel
  }

  const { body } = await request.GET<GetPipelineResponse>(
    `/api/v1/pipelines/${pipeline}`
  )
  return body.flow
}

export const savePipeline = async (pipeline: string, flow: PipelineModel): Promise<void> => {
  type SavePipelineRequest = {
    flow: PipelineModel
  }

  type SavePipelineResponse = {
    success: boolean
  }

  await request.PUT<SavePipelineRequest, SavePipelineResponse>(
    `/api/v1/pipelines/${pipeline}`,
    { flow }
  )
}

export const deletePipeline = async (pipeline: string): Promise<void> => {
  type DeletePipelineResponse = {
    success: boolean
  }

  await request.DELETE<DeletePipelineResponse>(`/api/v1/pipelines/${pipeline}`)
}
