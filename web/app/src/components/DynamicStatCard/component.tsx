import React from 'react'

import Skeleton from '@mui/material/Skeleton'

import AuthenticatedAwait from '@/components/AuthenticatedAwait/component'

import { DynamicStatCardProps } from './types'

import StatCard from '../StatCard/component'

const DynamicStatCard = <T,>({
  icon,
  title,
  href,
  resolver,
  renderer,
}: DynamicStatCardProps<T>) => (
  <React.Suspense
    fallback={
      <StatCard icon={icon} title={title} value={<Skeleton />} href={href} />
    }
  >
    <AuthenticatedAwait resolve={resolver()}>
      {(data) => (
        <StatCard
          icon={icon}
          title={title}
          value={renderer(data)}
          href={href}
        />
      )}
    </AuthenticatedAwait>
  </React.Suspense>
)

export default DynamicStatCard
