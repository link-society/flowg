import { useEffect } from 'react'
import { LoaderFunction, useLoaderData, useNavigate } from 'react-router'

import Backdrop from '@mui/material/Backdrop'
import CircularProgress from '@mui/material/CircularProgress'

import * as configApi from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

type LoaderData = {
  streams: string[]
}

export const loader: LoaderFunction = loginRequired(async () => {
  const streams = Object.keys(await configApi.listStreams())
  streams.sort((a, b) => a.localeCompare(b))
  return { streams }
})

const StreamSectionView = () => {
  const { streams } = useLoaderData() as LoaderData
  const navigate = useNavigate()

  useEffect(() => {
    if (streams.length > 0) {
      navigate(`/web/streams/${streams[0]}`, { replace: true })
    }
  }, [])

  return (
    <>
      {streams.length > 0 ? (
        <Backdrop open={true}>
          <CircularProgress color="inherit" />
        </Backdrop>
      ) : (
        <div className="w-full h-full flex flex-col items-center justify-center gap-5">
          <h1 className="text-3xl font-semibold">
            No stream found, send some logs.
          </h1>
        </div>
      )}
    </>
  )
}

export default StreamSectionView
