import { Outlet, useLoaderData, useLocation } from 'react-router'

import * as authApi from '@/lib/api/operations/auth'

import ProfileModel from '@/lib/models/ProfileModel'

import { loginRequired } from '@/lib/decorators/loaders'

import NavBar from '@/components/NavBar/component'
import PageFooter from '@/components/PageFooter/component'
import ProfileProvider from '@/components/ProfileProvider/component'

import { AppLayoutContainer } from './styles'

export const loader = async () => {
  return await loginRequired(authApi.whoami)()
}

const AppLayout = () => {
  const profile = useLoaderData() as ProfileModel
  const location = useLocation()

  return (
    <ProfileProvider value={profile}>
      <AppLayoutContainer>
        <NavBar />

        <main>
          <Outlet key={location.key} />
        </main>

        <PageFooter />
      </AppLayoutContainer>
    </ProfileProvider>
  )
}

export default AppLayout
