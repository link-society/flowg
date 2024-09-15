import { Outlet, redirect, useLoaderData } from 'react-router-dom'

import { NavBar } from '@/components/navbar'
import { ProfileProvider } from '@/lib/context/profile'

import { UnauthenticatedError } from '@/lib/api/errors'
import * as authApi from '@/lib/api/operations/auth'
import { ProfileModel } from '@/lib/models'

export const loader = async () => {
  try {
    return await authApi.whoami()
  }
  catch (error) {
    if (error instanceof UnauthenticatedError) {
      return redirect('/web/login')
    }
    else {
      throw error
    }
  }
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

        <footer className="p-3 bg-gray-300">
          <div className="text-center">
            <p>Footer</p>
          </div>
        </footer>
      </div>
    </ProfileProvider>
  )
}
