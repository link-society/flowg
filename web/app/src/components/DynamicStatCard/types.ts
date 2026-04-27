import { ReactNode } from 'react'

export type DynamicStatCardProps<T> = {
  icon: ReactNode
  title: ReactNode
  href: string
  resolver: () => Promise<T>
  renderer: (data: T) => ReactNode
}
