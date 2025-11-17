import { LoaderFunction } from 'react-router'

import * as configApi from '@/lib/api/operations/config'
import { loginRequired } from '@/lib/decorators/loaders'

export type LoaderData = {
  streams: string[]
  fields?: string[]
}

export const loader: LoaderFunction = loginRequired(async ({ params }) => {
  const streams = Object.keys(await configApi.listStreams())
  streams.sort((a, b) => a.localeCompare(b))

  if (params.stream !== undefined) {
    const fields = await configApi.listStreamFields(params.stream)
    return { streams, fields }
  }

  return { streams }
})
