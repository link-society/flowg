import { LoaderFunction } from 'react-router-dom'

import { loginRequired } from '@/lib/decorators/loaders'
import * as configApi from '@/lib/api/operations/config'

export type LoaderData = {
  transformers: string[]
  currentTransformer?: {
    name: string
    script: string
  }
}

export const loader: LoaderFunction = loginRequired(
  async ({ params }): Promise<LoaderData> => {
    const transformers = await configApi.listTransformers()

    if (params.transformer !== undefined) {
      if (!transformers.includes(params.transformer)) {
        throw new Response(
          `Transformer ${params.transformer} not found`,
          { status: 404 },
        )
      }

      const script = await configApi.getTransformer(params.transformer)
      return {
        transformers,
        currentTransformer: {
          name: params.transformer,
          script,
        },
      }
    }

    return { transformers }
  },
)
