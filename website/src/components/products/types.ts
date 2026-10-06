import type { IconName } from '@site/src/components/Icon/types'

export type ProductId = 'on-premise' | 'saas' | 'mcp'

export type ProductStep = Readonly<{ title: string; description: string }>
export type ProductQuestion = Readonly<{ question: string; answer: string }>
export type ProductClosing = Readonly<{ title: string; description: string }>

export type ProductFeature = Readonly<{
  icon: IconName
  title: string
  text: string
}>

export type Product = Readonly<{
  id: ProductId
  name: string
  shortName: string
  icon: IconName
  category: string
  title: string
  emphasis: string
  description: string
  action: string
  subject: string
  actionNote: string
  before: readonly string[]
  after: readonly string[]
  featuresTitle: string
  features: readonly ProductFeature[]
  offer: ProductClosing
  steps: readonly ProductStep[]
  faqs: readonly ProductQuestion[]
  closing: ProductClosing
}>
