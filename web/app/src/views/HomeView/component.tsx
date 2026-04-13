import { Typography } from '@mui/material'

import AccountTreeIcon from '@mui/icons-material/AccountTree'
import FilterAltIcon from '@mui/icons-material/FilterAlt'
import ForwardToInboxIcon from '@mui/icons-material/ForwardToInbox'
import StorageIcon from '@mui/icons-material/Storage'

import * as configApi from '@/lib/api/operations/config'

import { useProfile } from '@/lib/hooks/profile'

import StreamConfigModel from '@/lib/models/StreamConfigModel'

import DynamicStatCard from '@/components/DynamicStatCard/component'

import { HomeViewContainer, HomeViewPermissionsWrapper } from './styles'

const HomeView = () => {
  console.log('dev')
  const { permissions } = useProfile()

  return (
    <HomeViewContainer>
      <Typography variant="titleLg" component="h1">
        <span className="text-3xl">Welcome to FlowG</span>
        <img src="/web/assets/logo.png" alt="Logo" className="h-8" />
      </Typography>

      <HomeViewPermissionsWrapper>
        {permissions.can_view_streams && (
          <DynamicStatCard
            icon={<StorageIcon />}
            title="Streams"
            href="/web/streams"
            resolver={() => configApi.listStreams()}
            renderer={(streams: { [stream: string]: StreamConfigModel }) => {
              return Object.keys(streams).length
            }}
          />
        )}
        {permissions.can_view_transformers && (
          <DynamicStatCard
            icon={<FilterAltIcon />}
            title="Transformers"
            href="/web/transformers"
            resolver={() => configApi.listTransformers()}
            renderer={(transformers: string[]) => transformers.length}
          />
        )}
        {permissions.can_view_forwarders && (
          <DynamicStatCard
            icon={<ForwardToInboxIcon />}
            title="Forwarders"
            href="/web/forwarders"
            resolver={() => configApi.listForwarders()}
            renderer={(forwarders: string[]) => forwarders.length}
          />
        )}
        {permissions.can_view_pipelines && (
          <DynamicStatCard
            icon={<AccountTreeIcon />}
            title="Pipelines"
            href="/web/pipelines"
            resolver={() => configApi.listPipelines()}
            renderer={(pipelines: string[]) => pipelines.length}
          />
        )}
      </HomeViewPermissionsWrapper>
    </HomeViewContainer>
  )
}

export default HomeView
