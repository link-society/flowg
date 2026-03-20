import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs'
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider'

import { Outlet } from 'react-router'

import DialogsProvider from '@/components/DialogsProvider'
import NotificationsProvider from '@/components/NotificationsProvider'

import { StyledBaseLayout } from './styles'

const BaseLayout = () => {
  return (
    <StyledBaseLayout>
      <DialogsProvider>
        <NotificationsProvider>
          <LocalizationProvider dateAdapter={AdapterDayjs}>
            <Outlet />
          </LocalizationProvider>
        </NotificationsProvider>
      </DialogsProvider>
    </StyledBaseLayout>
  )
}

export default BaseLayout
