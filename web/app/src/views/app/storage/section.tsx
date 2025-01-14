import { useEffect } from 'react'
import { useLoaderData, useNavigate } from 'react-router'

import Backdrop from '@mui/material/Backdrop'
import CircularProgress from '@mui/material/CircularProgress'

import { LoaderData } from './loader'
import { NewStreamButton } from './new-btn'

export const StreamView = () => {
  const navigate = useNavigate()

  const { streams } = useLoaderData() as LoaderData
  const streamNames = Object.keys(streams)

  useEffect(
    () => {
      if (streamNames.length > 0) {
        navigate(`/web/storage/${streamNames[0]}`)
      }
    },
    [],
  )

  return (
    <>
      {streamNames.length > 0
        ? (
          <Backdrop open={true}>
            <CircularProgress color="inherit" />
          </Backdrop>
        )
        : (
          <div className="w-full h-full flex flex-col items-center justify-center gap-5">
            <h1 className="text-3xl font-semibold">No stream found, create one</h1>

            <NewStreamButton
              onStreamCreated={(name) => {
                window.location.pathname = `/web/storage/${name}`
              }}
            />
          </div>
        )
      }
    </>
  )
}
