import { LoaderFunction } from 'react-router'

import * as tokenApi from '@/lib/api/operations/token'
import { loginRequired } from '@/lib/decorators/loaders'

export type LoaderData = {
  tokens: string[]
}

export const loader: LoaderFunction = loginRequired(
  async (): Promise<LoaderData> => {
    const tokens = await tokenApi.listTokens()
    return { tokens }
  },
)
