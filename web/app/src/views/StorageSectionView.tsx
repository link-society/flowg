import { useEffect } from 'react'
import { LoaderFunction, useLoaderData, useNavigate } from 'react-router'

import Backdrop from '@mui/material/Backdrop'
import CircularProgress from '@mui/material/CircularProgress'

import * as configApi from '@/lib/api/operations/config'

import StreamConfigModel from '@/lib/models/StreamConfigModel'

import { loginRequired } from '@/lib/decorators/loaders'

import ButtonNewStreamConfig from '@/components/ButtonNewStreamConfig'

type LoaderData = {
  streams: Record<string, StreamConfigModel>
}

export const loader: LoaderFunction = loginRequired(async () => {
  const streams = await configApi.listStreams()
  return { streams }
})

const StorageSectionView = () => {
  const navigate = useNavigate()

  const { streams } = useLoaderData() as LoaderData
  const streamNames = Object.keys(streams)

  useEffect(() => {
    if (streamNames.length > 0) {
      navigate(`/web/storage/${streamNames[0]}`, { replace: true })
    }
  }, [])

  return (
    <>
      {streamNames.length > 0 ? (
        <Backdrop open={true}>
          <CircularProgress color="inherit" />
        </Backdrop>
      ) : (
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
      )}
    </>
  )
}

export default StorageSectionView
