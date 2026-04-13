import { Typography } from '@mui/material'

import { LoaderFunction, useLoaderData } from 'react-router'

import AccountTreeIcon from '@mui/icons-material/AccountTree'
import FilterAltIcon from '@mui/icons-material/FilterAlt'
import ForwardToInboxIcon from '@mui/icons-material/ForwardToInbox'
import StorageIcon from '@mui/icons-material/Storage'

import * as authApi from '@/lib/api/operations/auth'
import * as configApi from '@/lib/api/operations/config'

import StreamConfigModel from '@/lib/models/StreamConfigModel'

import { loginRequired } from '@/lib/decorators/loaders'

import StatCard from '@/components/StatCard/component'

import { HomeViewContainer, HomeViewPermissionsWrapper } from './styles'

type HomeViewData = {
  streams: { [stream: string]: StreamConfigModel } | null
  transformers: string[] | null
  forwarders: string[] | null
  pipelines: string[] | null
}

export const loader: LoaderFunction = loginRequired(async () => {
  const profile = await authApi.whoami()
  const { permissions } = profile

  const [streams, transformers, forwarders, pipelines] = await Promise.all([
    permissions.can_view_streams ? configApi.listStreams() : null,
    permissions.can_view_transformers ? configApi.listTransformers() : null,
    permissions.can_view_forwarders ? configApi.listForwarders() : null,
    permissions.can_view_pipelines ? configApi.listPipelines() : null,
  ])

  return { streams, transformers, forwarders, pipelines }
})

const HomeView = () => {
  const { streams, transformers, forwarders, pipelines } =
    useLoaderData<HomeViewData>()

  return (
    <HomeViewContainer variant="page">
      <Typography variant="titleLg" component="h1">
        <span className="text-3xl">Welcome to FlowG</span>
        <img src="/web/assets/logo.png" alt="Logo" className="h-8" />
      </Typography>

      <HomeViewPermissionsWrapper>
        {streams !== null && (
          <StatCard
            icon={<StorageIcon />}
            title="Streams"
            value={Object.keys(streams).length}
            to="/web/streams"
          />
        )}
        {transformers !== null && (
          <StatCard
            icon={<FilterAltIcon />}
            title="Transformers"
            value={transformers.length}
            to="/web/transformers"
          />
        )}
        {forwarders !== null && (
          <StatCard
            icon={<ForwardToInboxIcon />}
            title="Forwarders"
            value={forwarders.length}
            to="/web/forwarders"
          />
        )}
        {pipelines !== null && (
          <StatCard
            icon={<AccountTreeIcon />}
            title="Pipelines"
            value={pipelines.length}
            to="/web/pipelines"
          />
        )}
      </HomeViewPermissionsWrapper>
    </HomeViewContainer>
  )
}

export default HomeView
