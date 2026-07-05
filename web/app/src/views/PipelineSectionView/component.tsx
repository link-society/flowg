import { useTranslation } from 'react-i18next'
import { LoaderFunction, redirect, useNavigate } from 'react-router'

import Typography from '@mui/material/Typography'

import * as configApi from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

import ButtonNewPipeline from '@/components/ButtonNewPipeline/component'

import { buildUrl } from '@/router'

import { PipelineSectionViewRoot } from './styles'

export const loader: LoaderFunction = loginRequired(async () => {
  const pipelines = await configApi.listPipelines()
  if (pipelines.length > 0) {
    return redirect(buildUrl(`/pipelines/${pipelines[0]}`))
  }

  return null
})

const PipelineSectionView = () => {
  const { t } = useTranslation()
  const navigate = useNavigate()

  return (
    <PipelineSectionViewRoot>
      <Typography variant="titleLg" component="h1">
        {t('pages.pipelines.empty')}
      </Typography>

      <ButtonNewPipeline
        onPipelineCreated={(name) => {
          navigate(buildUrl(`/pipelines/${name}`))
        }}
      />
    </PipelineSectionViewRoot>
  )
}

export default PipelineSectionView
