import { LoaderFunction } from 'react-router'

import * as configApi from '@/lib/api/operations/config'
import { loginRequired } from '@/lib/decorators/loaders'
import { StreamConfigModel } from '@/lib/models/storage'

export type LoaderData = {
  streams: Record<string, StreamConfigModel>
  currentStream?: string
}

export const loader: LoaderFunction = loginRequired(async ({ params }) => {
  const streams = await configApi.listStreams()

  if (params.stream !== undefined) {
    if (streams[params.stream] === undefined) {
      throw new Response(`Stream ${params.stream} not found`, { status: 404 })
    }

    return {
      streams,
      currentStream: params.stream,
    }
  }

  return { streams }
})
