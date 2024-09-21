import { LoaderFunction } from 'react-router-dom'

import { loginRequired } from '@/lib/decorators/loaders'
import * as configApi from '@/lib/api/operations/config'

export type LoaderData = {
  streams: string[]
  fields?: string[]
}

export const loader: LoaderFunction = loginRequired(
  async ({ params }) => {
    const streams = Object.keys(await configApi.listStreams())
    streams.sort()

    if (params.stream !== undefined) {
      const fields = await configApi.listStreamFields(params.stream)
      return { streams, fields }
    }

    return { streams }
  },
)
