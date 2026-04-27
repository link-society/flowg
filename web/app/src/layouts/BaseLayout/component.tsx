import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs'
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider'

import { Outlet } from 'react-router'

import DialogsProvider from '@/components/DialogsProvider'
import NotificationsProvider from '@/components/NotificationsProvider'

import { BaseLayoutContainer } from './styles'

const BaseLayout = () => {
  return (
    <BaseLayoutContainer>
      <DialogsProvider>
        <NotificationsProvider>
          <LocalizationProvider dateAdapter={AdapterDayjs}>
            <Outlet />
          </LocalizationProvider>
        </NotificationsProvider>
      </DialogsProvider>
    </BaseLayoutContainer>
  )
}

export default BaseLayout
