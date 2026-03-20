import { Outlet, useLoaderData, useLocation } from 'react-router'

import * as authApi from '@/lib/api/operations/auth'

import ProfileModel from '@/lib/models/ProfileModel'

import { loginRequired } from '@/lib/decorators/loaders'

import NavBar from '@/components/NavBar'
import PageFooter from '@/components/PageFooter'
import ProfileProvider from '@/components/ProfileProvider'

import { StyledAppLayout } from './styles'

export const loader = async () => {
  return await loginRequired(authApi.whoami)()
}

const AppLayout = () => {
  const profile = useLoaderData() as ProfileModel
  const location = useLocation()

  return (
    <ProfileProvider value={profile}>
      <StyledAppLayout>
        <NavBar />

        <main>
          <Outlet key={location.key} />
        </main>

        <PageFooter />
      </StyledAppLayout>
    </ProfileProvider>
  )
}

export default AppLayout
