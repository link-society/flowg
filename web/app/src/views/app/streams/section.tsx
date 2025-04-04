import { useEffect } from 'react'
import { useLoaderData, useNavigate } from 'react-router'

import Backdrop from '@mui/material/Backdrop'
import CircularProgress from '@mui/material/CircularProgress'

import { LoaderData } from './loader'

export const StreamView = () => {
  const { streams } = useLoaderData() as LoaderData
  const navigate = useNavigate()

  useEffect(() => {
    if (streams.length > 0) {
      navigate(`/web/streams/${streams[0]}`)
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
