import type paths from './icons'

export type IconName = keyof typeof paths

export type IconProps = Readonly<{
  name: IconName
  size?: number
}>
