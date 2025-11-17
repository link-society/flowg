import { Outlet, useLoaderData } from 'react-router'

import * as authApi from '@/lib/api/operations/auth'

import ProfileModel from '@/lib/models/ProfileModel'

import { loginRequired } from '@/lib/decorators/loaders'

import NavBar from '@/components/NavBar'
import PageFooter from '@/components/PageFooter'
import ProfileProvider from '@/components/ProfileProvider'

export const loader = async () => {
  return await loginRequired(authApi.whoami)()
}

const AppLayout = () => {
  const profile = useLoaderData() as ProfileModel

  return (
    <ProfileProvider value={profile}>
      <div className="h-full flex flex-col overflow-hidden">
        <NavBar />

        <main className="grow shrink h-0">
          <Outlet />
        </main>

        <PageFooter />
      </div>
    </ProfileProvider>
  )
}

export default AppLayout
