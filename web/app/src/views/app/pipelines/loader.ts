import { LoaderFunction } from 'react-router-dom'

import { loginRequired } from '@/lib/decorators/loaders'
import * as configApi from '@/lib/api/operations/config'

export const loader: LoaderFunction = async ({ params }) => {
  const pipelines = await loginRequired(configApi.listPipelines)()

  if (params.pipeline !== undefined) {
    if (!pipelines.includes(params.pipeline)) {
      throw new Response(
        `Pipeline ${params.pipeline} not found`,
        { status: 404 },
      )
    }

    const flow = await loginRequired(configApi.getPipeline)(params.pipeline)
    return {
      pipelines,
      currentPipeline: {
        name: params.pipeline,
        flow,
      },
    }
  }

  return { pipelines }
}
