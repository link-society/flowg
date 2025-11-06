import { LoaderFunction } from 'react-router'

import * as configApi from '@/lib/api/operations/config'
import { loginRequired } from '@/lib/decorators/loaders'

export type LoaderData = {
  pipelines: string[]
}

export const loader: LoaderFunction = loginRequired(
  async (): Promise<LoaderData> => {
    const pipelines = await configApi.listPipelines()

    return { pipelines }
  }
)
