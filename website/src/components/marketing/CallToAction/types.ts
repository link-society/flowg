import type { ReactNode } from 'react'

export type CallToActionVariant = 'primary' | 'demo' | 'product' | 'text'

export type CallToActionProps = Readonly<{
  href: string
  children: ReactNode
  variant?: CallToActionVariant
  className?: string
}>
