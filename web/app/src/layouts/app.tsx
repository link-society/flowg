import { Outlet, useLoaderData } from 'react-router'

import { NavBar } from '@/components/navbar'
import { Footer } from '@/components/footer'
import { ProfileProvider } from '@/lib/context/profile'

import * as authApi from '@/lib/api/operations/auth'
import { loginRequired } from '@/lib/decorators/loaders'
import { ProfileModel } from '@/lib/models'

export const loader = async () => {
  return await loginRequired(authApi.whoami)()
}

export const AppLayout = () => {
  const profile = useLoaderData() as ProfileModel

  return (
    <ProfileProvider value={profile}>
      <div className="h-full flex flex-col overflow-hidden">
        <NavBar />

        <main className="flex-grow flex-shrink h-0">
          <Outlet />
        </main>

        <Footer />
      </div>
    </ProfileProvider>
  )
}
