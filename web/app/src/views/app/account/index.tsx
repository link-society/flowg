import { useLoaderData } from 'react-router'

import Grid from '@mui/material/Grid'

import { ProfileInfo } from './profile-info'
import { TokenList } from './token-list'

import { LoaderData } from './loader'

export const AccountView = () => {
  const { tokens } = useLoaderData() as LoaderData

  return (
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
  )
}
