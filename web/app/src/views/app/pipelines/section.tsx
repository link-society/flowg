import { useEffect } from 'react'
import { useLoaderData, useNavigate } from 'react-router'

import Backdrop from '@mui/material/Backdrop'
import CircularProgress from '@mui/material/CircularProgress'

import { LoaderData } from './loader'
import { NewPipelineButton } from './new-btn'

export const PipelineView = () => {
  const { pipelines } = useLoaderData() as LoaderData
  const navigate = useNavigate()

  useEffect(() => {
    if (pipelines.length > 0) {
      navigate(`/web/pipelines/${pipelines[0]}`, { replace: true })
    }
  }, [])

  return (
    <>
      {pipelines.length > 0 ? (
        <Backdrop open={true}>
          <CircularProgress color="inherit" />
        </Backdrop>
      ) : (
        <div className="w-full h-full flex flex-col items-center justify-center gap-5">
          <h1 className="text-3xl font-semibold">
            No pipeline found, create one
          </h1>

          <NewPipelineButton
            onPipelineCreated={(name) => {
              globalThis.location.pathname = `/web/pipelines/${name}`
            }}
          />
        </div>
      )}
    </>
  )
}
