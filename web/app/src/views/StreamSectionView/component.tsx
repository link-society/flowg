import { Typography } from '@mui/material'

import { LoaderFunction, redirect } from 'react-router'

import * as configApi from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

import { StreamSectionViewContainer } from './styles'

export const loader: LoaderFunction = loginRequired(async () => {
  const streams = Object.keys(await configApi.listStreams())
  streams.sort((a, b) => a.localeCompare(b))
  if (streams.length > 0) {
    return redirect(`/web/streams/${streams[0]}`)
  }

  return null
})

const StreamSectionView = () => {
  return (
    <StreamSectionViewContainer variant="page">
      <Typography variant="titleLg" component="h1">
        No stream found, send some logs.
      </Typography>
    </StreamSectionViewContainer>
  )
}

export default StreamSectionView
