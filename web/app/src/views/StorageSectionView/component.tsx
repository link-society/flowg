import { useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { LoaderFunction, useLoaderData, useNavigate } from 'react-router'

import Typography from '@mui/material/Typography'

import * as configApi from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

import ButtonNewStreamConfig from '@/components/ButtonNewStreamConfig/component'

import { buildUrl } from '@/router'

import { StorageSectionViewRoot } from './styles'

export const loader: LoaderFunction = loginRequired(async () => {
  const streams = await configApi.listStreams()
  const streamNames = Object.keys(streams)
  if (streamNames.length > 0) {
    return { redirectTo: buildUrl(`/storage/${streamNames[0]}`) }
  }

  return { redirectTo: null }
})

const StorageSectionView = () => {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { redirectTo } = useLoaderData<{ redirectTo: string | null }>()

  useEffect(() => {
    if (redirectTo !== null) {
      navigate(redirectTo, { replace: true })
    }
  }, [redirectTo])

  return (
    <StorageSectionViewRoot>
      <Typography variant="titleLg" component="h1">
        {t('pages.storage.empty')}
      </Typography>

      <ButtonNewStreamConfig
        onStreamConfigCreated={(name) => {
          navigate(buildUrl(`/storage/${name}`))
        }}
      />
    </StorageSectionViewRoot>
  )
}

export default StorageSectionView
