import { Outlet, useLoaderData } from 'react-router'

import * as authApi from '@/lib/api/operations/auth'
import { ProfileProvider } from '@/lib/context/profile'
import { loginRequired } from '@/lib/decorators/loaders'
import { ProfileModel } from '@/lib/models/auth'

import { Footer } from '@/components/footer'
import { NavBar } from '@/components/navbar'

export const loader = async () => {
  return await loginRequired(authApi.whoami)()
}

export const AppLayout = () => {
  const profile = useLoaderData() as ProfileModel

  return (
    <ProfileProvider value={profile}>
      <div className="h-full flex flex-col overflow-hidden">
        <NavBar />

        <main className="grow shrink h-0">
          <Outlet />
        </main>

        <Footer />
      </div>
    </ProfileProvider>
  )
}
