import { useEffect } from 'react'
import { useLoaderData, useNavigate } from 'react-router-dom'

import Backdrop from '@mui/material/Backdrop'
import CircularProgress from '@mui/material/CircularProgress'

import { LoaderData } from './loader'
import { NewAlertButton } from './new-btn'

export const AlertView = () => {
  const { alerts } = useLoaderData() as LoaderData
  const navigate = useNavigate()

  useEffect(
    () => {
      if (alerts.length > 0) {
        navigate(`/web/alerts/${alerts[0]}`)
      }
    },
    [],
  )

  return (
    <>
      {alerts.length > 0
        ? (
          <Backdrop open={true}>
            <CircularProgress color="inherit" />
          </Backdrop>
        )
        : (
          <div className="w-full h-full flex flex-col items-center justify-center gap-5">
            <h1 className="text-3xl font-semibold">No alert found, create one</h1>

            <NewAlertButton
              onAlertCreated={(name) => {
                window.location.pathname = `/web/alerts/${name}`
              }}
            />
          </div>
        )
      }
    </>
  )
}
