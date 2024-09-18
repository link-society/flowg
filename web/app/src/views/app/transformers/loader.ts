import { LoaderFunction } from 'react-router-dom'

import { loginRequired } from '@/lib/decorators/loaders'
import * as configApi from '@/lib/api/operations/config'

export const loader: LoaderFunction = async ({ params }) => {
  const transformers = await loginRequired(configApi.listTransformers)()

  if (params.transformer !== undefined) {
    if (!transformers.includes(params.transformer)) {
      throw new Response(
        `Transformer ${params.transformer} not found`,
        { status: 404 },
      )
    }

    const script = await loginRequired(configApi.getTransformer)(params.transformer)
    return {
      transformers,
      currentTransformer: {
        name: params.transformer,
        script,
      },
    }
  }

  return { transformers }
}
