import { LoaderFunction, useNavigate, redirect } from 'react-router'

import * as configApi from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

import ButtonNewPipeline from '@/components/ButtonNewPipeline'

export const loader: LoaderFunction = loginRequired(async () => {
  const pipelines = await configApi.listPipelines()
  if (pipelines.length > 0) {
    throw redirect(`/web/pipelines/${pipelines[0]}`)
  }
})

const PipelineSectionView = () => {
  const navigate = useNavigate()

  return (
    <div className="w-full h-full flex flex-col items-center justify-center gap-5">
      <h1 className="text-3xl font-semibold">
        No pipeline found, create one
      </h1>

      <ButtonNewPipeline
        onPipelineCreated={(name) => {
          navigate(`/web/pipelines/${name}`)
        }}
      />
    </div>
  )
}

export default PipelineSectionView
