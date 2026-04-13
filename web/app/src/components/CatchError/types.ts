import { ReactNode } from 'react'

export type CatchErrorProps<E extends Error> = {
  errorType: new (...args: any[]) => E
  fallback: ReactNode
}
