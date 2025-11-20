import { LoaderFunction, redirect, useNavigate } from 'react-router'

import * as configApi from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

import ButtonNewStreamConfig from '@/components/ButtonNewStreamConfig'

export const loader: LoaderFunction = loginRequired(async () => {
  const streams = await configApi.listStreams()
  const streamNames = Object.keys(streams)
  if (streamNames.length > 0) {
    throw redirect(`/web/storage/${streamNames[0]}`)
  }
})

const StorageSectionView = () => {
  const navigate = useNavigate()

  return (
    <div className="w-full h-full flex flex-col items-center justify-center gap-5">
      <h1 className="text-3xl font-semibold">
        No stream found, create one
      </h1>

      <ButtonNewStreamConfig
        onStreamConfigCreated={(name) => {
          navigate(`/web/storage/${name}`)
        }}
      />
    </div>
  )
}

export default StorageSectionView
