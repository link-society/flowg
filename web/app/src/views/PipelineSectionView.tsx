import { useEffect } from 'react'
import { LoaderFunction, useLoaderData, useNavigate } from 'react-router'

import Backdrop from '@mui/material/Backdrop'
import CircularProgress from '@mui/material/CircularProgress'

import * as configApi from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

import ButtonNewPipeline from '@/components/ButtonNewPipeline'

type LoaderData = {
  pipelines: string[]
}

export const loader: LoaderFunction = loginRequired(
  async (): Promise<LoaderData> => {
    const pipelines = await configApi.listPipelines()
    return { pipelines }
  }
)

const PipelineSectionView = () => {
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

          <ButtonNewPipeline
            onPipelineCreated={(name) => {
              globalThis.location.pathname = `/web/pipelines/${name}`
            }}
          />
        </div>
      )}
    </>
  )
}

export default PipelineSectionView
