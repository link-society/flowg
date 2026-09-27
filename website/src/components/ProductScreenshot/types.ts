import type { IconName } from '@site/src/components/Icon/types'

export type ProductScreenshotProps = Readonly<{
  icon: IconName
  title: string
  src: string
  alt: string
  caption: string
  width?: number
  height?: number
  reverse?: boolean
}>
