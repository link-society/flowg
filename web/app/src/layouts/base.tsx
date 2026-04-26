import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs'
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider'

import { Outlet } from 'react-router'

import { NotificationsProvider } from '@toolpad/core/useNotifications'

import DialogsProvider from '@/components/DialogsProvider'

const BaseLayout = () => {
  return (
    <div className="h-full flex flex-col overflow-hidden">
      <DialogsProvider>
        <NotificationsProvider>
          <LocalizationProvider dateAdapter={AdapterDayjs}>
            <Outlet />
          </LocalizationProvider>
        </NotificationsProvider>
      </DialogsProvider>
    </div>
  )
}

export default BaseLayout
