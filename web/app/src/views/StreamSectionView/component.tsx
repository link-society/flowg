import { Typography } from '@mui/material'

import { useEffect } from 'react'
import { LoaderFunction, useLoaderData, useNavigate } from 'react-router'

import * as configApi from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

import { StreamSectionViewContainer } from './styles'

export const loader: LoaderFunction = loginRequired(async () => {
  const streams = Object.keys(await configApi.listStreams())
  streams.sort((a, b) => a.localeCompare(b))
  if (streams.length > 0) {
    return { redirectTo: `/web/streams/${streams[0]}` }
  }

  return { redirectTo: null }
})

const StreamSectionView = () => {
  const navigate = useNavigate()
  const { redirectTo } = useLoaderData<{ redirectTo: string | null }>()

  useEffect(() => {
    if (redirectTo !== null) {
      navigate(redirectTo, { replace: true })
    }
  }, [redirectTo])

  return (
    <StreamSectionViewContainer>
      <Typography variant="titleLg" fontWeight={700} component="h1">
        No stream found, send some logs.
      </Typography>
    </StreamSectionViewContainer>
  )
}

export default StreamSectionView
