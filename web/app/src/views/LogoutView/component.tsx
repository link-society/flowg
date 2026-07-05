import { useTranslation } from 'react-i18next'
import { redirect } from 'react-router'

import Typography from '@mui/material/Typography'

import * as authApi from '@/lib/api/operations/auth'

import { buildUrl } from '@/router'

import { LogoutViewRoot } from './styles'

export const loader = async () => {
  await authApi.logout()
  throw redirect(buildUrl('/login'))
}

const LogoutView = () => {
  const { t } = useTranslation()

  return (
    <LogoutViewRoot>
      <Typography variant="text">{t('pages.logout.message')}</Typography>
    </LogoutViewRoot>
  )
}

export default LogoutView
