import { redirect } from 'react-router'

import Typography from '@mui/material/Typography'

import * as authApi from '@/lib/api/operations/auth'

import { LogoutViewRoot } from './styles'

export const loader = async () => {
  await authApi.logout()
  throw redirect('/web/login')
}

const LogoutView = () => {
  return (
    <LogoutViewRoot>
      <Typography variant="text">You are being logged out...</Typography>
    </LogoutViewRoot>
  )
}

export default LogoutView
