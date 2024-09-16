import React, { ReactElement, ReactNode } from 'react'

import Skeleton from '@mui/material/Skeleton'

import { AuthenticatedAwait } from '@/components/routing/await'

import { StatCard } from './statcard'

type DynamicStatCardProps<T> = {
  icon: ReactNode
  title: ReactNode
  href: string
  resolver: () => Promise<T>
  renderer: (data: T) => ReactNode
}

type C = <T>(props: DynamicStatCardProps<T>) => ReactElement
export const DynamicStatCard: C = (props) => (
  <React.Suspense
    fallback={
      <StatCard
        icon={props.icon}
        title={props.title}
        value={<Skeleton variant="text" />}
        href={props.href}
      />
    }
  >
    <AuthenticatedAwait resolve={props.resolver()}>
      {(data) => (
        <StatCard
          icon={props.icon}
          title={props.title}
          value={props.renderer(data)}
          href={props.href}
        />
      )}
    </AuthenticatedAwait>
  </React.Suspense>
)
