import type { ReactNode } from 'react'

export type ProductFeatureProps = Readonly<{
  eyebrow: string
  title: ReactNode
  children: ReactNode
  screenshot: ReactNode
  reverse?: boolean
}>
