import { Typography } from '@mui/material'

import { LoaderFunction, useLoaderData } from 'react-router'

import AccountTreeIcon from '@mui/icons-material/AccountTree'
import FilterAltIcon from '@mui/icons-material/FilterAlt'
import ForwardToInboxIcon from '@mui/icons-material/ForwardToInbox'
import StorageIcon from '@mui/icons-material/Storage'

import * as authApi from '@/lib/api/operations/auth'
import * as configApi from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

import StatCard from '@/components/StatCard/component'

import { buildUrl } from '@/router'

import { HomeViewContainer, HomeViewPermissionsWrapper } from './styles'
import { HomeViewData } from './types'

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
        <span>Welcome to FlowG</span>
        <img src={buildUrl('/assets/logo.png')} alt="Logo FlowG" />
      </Typography>

      <HomeViewPermissionsWrapper>
        {streams !== null && (
          <StatCard
            icon={<StorageIcon />}
            title="Streams"
            value={Object.keys(streams).length}
            to={buildUrl('/streams')}
          />
        )}
        {transformers !== null && (
          <StatCard
            icon={<FilterAltIcon />}
            title="Transformers"
            value={transformers.length}
            to={buildUrl('/transformers')}
          />
        )}
        {forwarders !== null && (
          <StatCard
            icon={<ForwardToInboxIcon />}
            title="Forwarders"
            value={forwarders.length}
            to={buildUrl('/forwarders')}
          />
        )}
        {pipelines !== null && (
          <StatCard
            icon={<AccountTreeIcon />}
            title="Pipelines"
            value={pipelines.length}
            to={buildUrl('/pipelines')}
          />
        )}
      </HomeViewPermissionsWrapper>
    </HomeViewContainer>
  )
}

export default HomeView
