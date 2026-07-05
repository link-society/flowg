import { Typography } from '@mui/material'

import { useTranslation } from 'react-i18next'
import { LoaderFunction, redirect } from 'react-router'

import * as configApi from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

import { buildUrl } from '@/router'

import { StreamSectionViewContainer } from './styles'

export const loader: LoaderFunction = loginRequired(async () => {
  const streams = Object.keys(await configApi.listStreams())
  streams.sort((a, b) => a.localeCompare(b))
  if (streams.length > 0) {
    return redirect(buildUrl(`/streams/${streams[0]}`))
  }

  return null
})

const StreamSectionView = () => {
  const { t } = useTranslation()

  return (
    <StreamSectionViewContainer variant="page">
      <Typography variant="titleLg" component="h1">
        {t('pages.streams.empty')}
      </Typography>
    </StreamSectionViewContainer>
  )
}

export default StreamSectionView
