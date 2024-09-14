import { Outlet, redirect } from 'react-router-dom'

import * as api from '@/lib/api'
import * as authApi from '@/lib/api/auth'

export const loader = async () => {
  try {
    const user = await authApi.whoami()
    return user
  }
  catch (error) {
    if (error instanceof api.UnauthenticatedError) {
      return redirect('/web/login')
    }
    else {
      throw error
    }
  }
}

export default function AppLayout() {
  return (
    <>
      <div>App</div>
      <Outlet />
    </>
  )
}
