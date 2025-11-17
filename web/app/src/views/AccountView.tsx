import { LoaderFunction, useLoaderData } from 'react-router'

import Grid from '@mui/material/Grid'

import * as tokenApi from '@/lib/api/operations/token'

import { loginRequired } from '@/lib/decorators/loaders'

import ProfileInfo from '@/components/ProfileInfo'
import TokenTable from '@/components/TokenTable'

type LoaderData = {
  tokens: string[]
}

export const loader: LoaderFunction = loginRequired(
  async (): Promise<LoaderData> => {
    const tokens = await tokenApi.listTokens()
    return { tokens }
  }
)

const AccountView = () => {
  const { tokens } = useLoaderData() as LoaderData

  return (
    <Grid container spacing={2} className="p-3 h-full max-lg:overflow-auto">
      <Grid size={{ xs: 12, md: 6 }} className="lg:h-full">
        <ProfileInfo />
      </Grid>
      <Grid size={{ xs: 12, md: 6 }} className="lg:h-full">
        <TokenTable tokens={tokens} />
      </Grid>
    </Grid>
  )
}

export default AccountView
