import type { ReactNode } from 'react'

import type { IconName } from '@site/src/components/Icon/types'

export type ProductCardProps = Readonly<{
  id: string
  name: string
  icon: IconName
  label: string
  headline: ReactNode
  description: string
  benefits: readonly { icon: IconName; text: string }[]
  servicesLabel: string
  services: readonly string[]
  action: { href: string; label: string }
  note: string
}>
