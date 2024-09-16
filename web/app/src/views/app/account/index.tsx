import { useLoaderData } from 'react-router-dom'

import Grid from '@mui/material/Grid2'

import { ProfileInfo } from './profile-info'
import { TokenList } from './token-list'

import * as tokenApi from '@/lib/api/operations/token'
import { loginRequired } from '@/lib/decorators/loaders'

export const loader = async () => {
  return await loginRequired(tokenApi.listTokens)()
}

export const AccountView = () => {
  const tokens = useLoaderData() as string[]

  return (
    <>
      <Grid
        container
        spacing={2}
        className="p-3 h-full max-lg:overflow-auto"
      >
        <Grid size={{ xs: 12, md: 6 }} className="lg:h-full">
          <ProfileInfo />
        </Grid>
        <Grid size={{ xs: 12, md: 6 }} className="lg:h-full">
          <TokenList tokens={tokens} />
        </Grid>
      </Grid>
    </>
  )
}
