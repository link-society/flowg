import { LoaderFunction, useLoaderData } from 'react-router'

import * as tokenApi from '@/lib/api/operations/token'

import { loginRequired } from '@/lib/decorators/loaders'

import ProfileInfo from '@/components/ProfileInfo/component'
import TokenTable from '@/components/TokenTable/component'

import { AccountViewContainer, AccountViewPanel } from './styles'
import { LoaderData } from './types'

export const loader: LoaderFunction = loginRequired(
  async (): Promise<LoaderData> => {
    const tokens = await tokenApi.listTokens()
    return { tokens }
  }
)

const AccountView = () => {
  const { tokens } = useLoaderData() as LoaderData

  return (
    <AccountViewContainer variant="page">
      <AccountViewPanel>
        <ProfileInfo />
      </AccountViewPanel>
      <AccountViewPanel>
        <TokenTable tokens={tokens} />
      </AccountViewPanel>
    </AccountViewContainer>
  )
}

export default AccountView
