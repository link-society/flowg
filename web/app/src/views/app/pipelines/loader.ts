import { LoaderFunction } from 'react-router'

import { loginRequired } from '@/lib/decorators/loaders'
import * as configApi from '@/lib/api/operations/config'
import { PipelineModel } from '@/lib/models/pipeline'

export type LoaderData = {
  pipelines: string[]
  currentPipeline?: {
    name: string
    flow: PipelineModel
  }
}

export const loader: LoaderFunction = loginRequired(
  async ({ params }): Promise<LoaderData> => {
    const pipelines = await configApi.listPipelines()

    if (params.pipeline !== undefined) {
      if (!pipelines.includes(params.pipeline)) {
        throw new Response(
          `Pipeline ${params.pipeline} not found`,
          { status: 404 },
        )
      }

      const flow = await configApi.getPipeline(params.pipeline)
      return {
        pipelines,
        currentPipeline: {
          name: params.pipeline,
          flow,
        },
      }
    }

    return { pipelines }
  },
)
