import { useEffect } from 'react'
import { LoaderFunction, useLoaderData, useNavigate } from 'react-router'

import Typography from '@mui/material/Typography'

import * as configApi from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

import ButtonNewStreamConfig from '@/components/ButtonNewStreamConfig/component'

import { StorageSectionViewRoot } from './styles'

export const loader: LoaderFunction = loginRequired(async () => {
  const streams = await configApi.listStreams()
  const streamNames = Object.keys(streams)
  if (streamNames.length > 0) {
    return { redirectTo: `/web/storage/${streamNames[0]}` }
  }

  return { redirectTo: null }
})

const StorageSectionView = () => {
  const navigate = useNavigate()
  const { redirectTo } = useLoaderData<{ redirectTo: string | null }>()

  useEffect(() => {
    if (redirectTo !== null) {
      navigate(redirectTo, { replace: true })
    }
  }, [redirectTo])

  return (
    <StorageSectionViewRoot>
      <Typography variant="titleLg" fontWeight={700}>
        No stream found, create one
      </Typography>

      <ButtonNewStreamConfig
        onStreamConfigCreated={(name) => {
          navigate(`/web/storage/${name}`)
        }}
      />
    </StorageSectionViewRoot>
  )
}

export default StorageSectionView
