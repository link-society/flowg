import { redirect } from 'react-router-dom'

import * as authApi from '@/lib/api/auth'

export const loader = async () => {
  await authApi.logout()
  return redirect('/web/login')
}

export default function LogoutView() {
  return (
    <div className="py-6">
      You are being logged out...
    </div>
  )
}