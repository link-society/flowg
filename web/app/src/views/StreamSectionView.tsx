import { LoaderFunction, redirect } from 'react-router'

import * as configApi from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

export const loader: LoaderFunction = loginRequired(async () => {
  const streams = Object.keys(await configApi.listStreams())
  streams.sort((a, b) => a.localeCompare(b))
  if (streams.length > 0) {
    throw redirect(`/web/streams/${streams[0]}`)
  }
})

const StreamSectionView = () => (
  <div className="w-full h-full flex flex-col items-center justify-center gap-5">
    <h1 className="text-3xl font-semibold">
      No stream found, send some logs.
    </h1>
  </div>
)

export default StreamSectionView
