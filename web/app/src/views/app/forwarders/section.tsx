import { useEffect } from 'react'
import { useLoaderData, useNavigate } from 'react-router'

import Backdrop from '@mui/material/Backdrop'
import CircularProgress from '@mui/material/CircularProgress'

import { LoaderData } from './loader'
import { NewForwarderButton } from './new-btn'

export const ForwarderView = () => {
  const { forwarders } = useLoaderData() as LoaderData
  const navigate = useNavigate()

  useEffect(
    () => {
      if (forwarders.length > 0) {
        navigate(`/web/forwarders/${forwarders[0]}`)
      }
    },
    [],
  )

  return (
    <>
      {forwarders.length > 0
        ? (
          <Backdrop open={true}>
            <CircularProgress color="inherit" />
          </Backdrop>
        )
        : (
          <div className="w-full h-full flex flex-col items-center justify-center gap-5">
            <h1 className="text-3xl font-semibold">No forwarder found, create one</h1>

            <NewForwarderButton
              onForwarderCreated={(name) => {
                navigate(`/web/forwarders/${name}`)
              }}
            />
          </div>
        )
      }
    </>
  )
}
