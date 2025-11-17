import { redirect } from 'react-router'

import * as authApi from '@/lib/api/operations/auth'

export const loader = async () => {
  await authApi.logout()
  return redirect('/web/login')
}

const LogoutView = () => {
  return <div className="py-6">You are being logged out...</div>
}

export default LogoutView
