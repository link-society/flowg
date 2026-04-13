import { useEffect } from 'react'
import { LoaderFunction, useLoaderData, useNavigate } from 'react-router'

import Typography from '@mui/material/Typography'

import * as configApi from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

import ButtonNewPipeline from '@/components/ButtonNewPipeline/component'

import { PipelineSectionViewRoot } from './styles'

export const loader: LoaderFunction = loginRequired(async () => {
  const pipelines = await configApi.listPipelines()
  if (pipelines.length > 0) {
    return { redirectTo: `/web/pipelines/${pipelines[0]}` }
  }

  return { redirectTo: null }
})

const PipelineSectionView = () => {
  const navigate = useNavigate()
  const { redirectTo } = useLoaderData<{ redirectTo: string | null }>()

  useEffect(() => {
    if (redirectTo !== null) {
      navigate(redirectTo, { replace: true })
    }
  }, [redirectTo])

  return (
    <PipelineSectionViewRoot>
      <Typography variant="titleLg" fontWeight={700} component="h1">
        No pipeline found, create one
      </Typography>

      <ButtonNewPipeline
        onPipelineCreated={(name) => {
          navigate(`/web/pipelines/${name}`)
        }}
      />
    </PipelineSectionViewRoot>
  )
}

export default PipelineSectionView
