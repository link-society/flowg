import * as request from '@/lib/api/request'

import ForwarderModel from '@/lib/models/ForwarderModel'
import PipelineModel from '@/lib/models/PipelineModel'
import { PipelineTrace } from '@/lib/models/PipelineTrace.ts'
import StreamConfigModel from '@/lib/models/StreamConfigModel'
import SystemConfigurationModel from '@/lib/models/SystemConfigurationModel.ts'

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

export const saveTransformer = async (
  transformer: string,
  script: string
): Promise<void> => {
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

export const listStreams = async (): Promise<{
  [stream: string]: StreamConfigModel
}> => {
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

export const getStreamConfig = async (
  stream: string
): Promise<StreamConfigModel> => {
  type GetStreamConfigResponse = {
    success: boolean
    config: StreamConfigModel
  }

  const { body } = await request.GET<GetStreamConfigResponse>({
    path: `/api/v1/streams/${stream}`,
  })

  return body.config
}

export const configureStream = async (
  stream: string,
  config: StreamConfigModel
): Promise<void> => {
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

export const listForwarders = async (): Promise<string[]> => {
  type ListForwardersResponse = {
    success: boolean
    forwarders: string[]
  }

  const { body } = await request.GET<ListForwardersResponse>({
    path: '/api/v1/forwarders',
  })
  return body.forwarders
}

export const getForwarder = async (
  forwarder: string
): Promise<ForwarderModel> => {
  type GetForwarderResponse = {
    success: boolean
    forwarder: ForwarderModel
  }

  const { body } = await request.GET<GetForwarderResponse>({
    path: `/api/v1/forwarders/${forwarder}`,
  })

  return body.forwarder
}

export const saveForwarder = async (
  name: string,
  forwarder: ForwarderModel
): Promise<void> => {
  type SaveForwarderRequest = {
    forwarder: ForwarderModel
  }

  type SaveForwarderResponse = {
    success: boolean
  }

  await request.PUT<SaveForwarderRequest, SaveForwarderResponse>({
    path: `/api/v1/forwarders/${name}`,
    body: { forwarder },
  })
}

export const deleteForwarder = async (forwarder: string): Promise<void> => {
  type DeleteForwarderResponse = {
    success: boolean
  }

  await request.DELETE<DeleteForwarderResponse>({
    path: `/api/v1/forwarders/${forwarder}`,
  })
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

export const savePipeline = async (
  pipeline: string,
  flow: PipelineModel
): Promise<void> => {
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

type TestPipelineResult =
  | {
      success: boolean
      trace: PipelineTrace
    }
  | { success: false; error: string }

export const testPipeline = async (
  pipeline: string,
  records: Record<string, string>[]
): Promise<TestPipelineResult> => {
  type TestPipelineRequest = {
    records: Record<string, string>[]
  }

  const { body } = await request.POST<TestPipelineRequest, TestPipelineResult>({
    path: `/api/v1/test/pipeline/${pipeline}`,
    body: { records },
  })

  return body
}

export const getSystemConfiguration =
  async (): Promise<SystemConfigurationModel> => {
    type getSystemConfigurationResponse = {
      success: boolean
      configuration: SystemConfigurationModel
    }

    const { body } = await request.GET<getSystemConfigurationResponse>({
      path: `/api/v1/system-configuration`,
    })
    return body.configuration
  }

export const saveSystemConfiguration = async (
  config: SystemConfigurationModel
): Promise<void> => {
  type saveSystemConfigurationResponse = {
    success: boolean
  }

  await request.PUT<SystemConfigurationModel, saveSystemConfigurationResponse>({
    path: `/api/v1/system-configuration`,
    body: config,
  })
}
