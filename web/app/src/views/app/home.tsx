import React, { ReactNode } from 'react'
import { Await } from 'react-router-dom'

import StorageIcon from '@mui/icons-material/Storage'
import FilterAltIcon from '@mui/icons-material/FilterAlt'
import NotificationsActiveIcon from '@mui/icons-material/NotificationsActive'
import AccountTreeIcon from '@mui/icons-material/AccountTree'

import Grid from '@mui/material/Grid2'
import Card from '@mui/material/Card'
import CardHeader from '@mui/material/CardHeader'
import CardContent from '@mui/material/CardContent'
import CardActions from '@mui/material/CardActions'
import Button from '@mui/material/Button'
import Skeleton from '@mui/material/Skeleton'

import { useProfile } from '@/lib/context/profile'

import * as configApi from '@/lib/api/operations/config'
import { StreamConfigModel } from '@/lib/models'

type StatCardProps = {
  icon: ReactNode
  title: ReactNode
  value: ReactNode
  href: string
}

const StatCard = ({ icon, title, value, href }: StatCardProps) => (
  <Card>
    <CardHeader
      title={
        <div
          className="
            flex items-center justify-center gap-3
            text-2xl font-semibold
          "
        >
          {icon}
          {title}
        </div>
      }
    />
    <CardContent className="!p-0 text-center text-3xl font-bold">
      {value}
      <hr className="mt-3" />
    </CardContent>
    <CardActions>
      <Button href={href} className="w-full">
        View More
      </Button>
    </CardActions>
  </Card>
)

export const HomeView = () => {
  const { permissions } = useProfile()

  return (
    <Grid container spacing={2} className="justify-center py-6">
      {permissions.can_view_streams && (
        <Grid size={{ sm: 12, md: 2 }}>
          <React.Suspense
            fallback={
              <StatCard
                icon={<StorageIcon />}
                title="Streams"
                value={<Skeleton variant="text" />}
                href="/web/streams"
              />
            }
          >
            <Await resolve={configApi.listStreams()}>
              {(streams: { [stream: string]: StreamConfigModel }) => (
                <StatCard
                  icon={<StorageIcon />}
                  title="Streams"
                  value={Object.keys(streams).length}
                  href="/web/streams"
                />
              )}
            </Await>
          </React.Suspense>
        </Grid>
      )}
      {permissions.can_view_transformers && (
        <Grid size={{ sm: 12, md: 2 }}>
          <React.Suspense
            fallback={
              <StatCard
                icon={<FilterAltIcon />}
                title="Transformers"
                value={<Skeleton variant="text" />}
                href="/web/transformers"
              />
            }
          >
            <Await resolve={configApi.listTransformers()}>
              {(transformers: string[]) => (
                <StatCard
                  icon={<FilterAltIcon />}
                  title="Transformers"
                  value={transformers.length}
                  href="/web/transformers"
                />
              )}
            </Await>
          </React.Suspense>
        </Grid>
      )}
      {permissions.can_view_alerts && (
        <Grid size={{ sm: 12, md: 2 }}>
          <React.Suspense
            fallback={
              <StatCard
                icon={<NotificationsActiveIcon />}
                title="Alerts"
                value={<Skeleton variant="text" />}
                href="/web/alerts"
              />
            }
          >
            <Await resolve={configApi.listAlerts()}>
              {(alerts: string[]) => (
                <StatCard
                  icon={<NotificationsActiveIcon />}
                  title="Alerts"
                  value={alerts.length}
                  href="/web/alerts"
                />
              )}
            </Await>
          </React.Suspense>
        </Grid>
      )}
      {permissions.can_view_pipelines && (
        <Grid size={{ sm: 12, md: 2 }}>
          <React.Suspense
            fallback={
              <StatCard
                icon={<AccountTreeIcon />}
                title="Pipelines"
                value={<Skeleton variant="text" />}
                href="/web/pipeliens"
              />
            }
          >
            <Await resolve={configApi.listPipelines()}>
              {(pipelines: string[]) => (
                <StatCard
                  icon={<AccountTreeIcon />}
                  title="Pipelines"
                  value={pipelines.length}
                  href="/web/pipeliens"
                />
              )}
            </Await>
          </React.Suspense>
        </Grid>
      )}
    </Grid>
  )
}
