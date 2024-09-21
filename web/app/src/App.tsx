import { RouterProvider } from 'react-router-dom'

import { createTheme, ThemeProvider } from '@mui/material/styles'
import * as colors from '@mui/material/colors'

import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs'
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider'

import { DialogsProvider } from '@toolpad/core/useDialogs'
import { NotificationsProvider } from '@toolpad/core/useNotifications'

import router from '@/router'

const theme = createTheme({
  shape: {
    borderRadius: 0,
  },
  palette: {
    primary: {
      main: colors.blue[800],
    },
    secondary: {
      main: colors.teal[400],

    }
  }
})

export default function App() {
  return (
    <ThemeProvider theme={theme}>
      <DialogsProvider>
        <NotificationsProvider>
          <LocalizationProvider dateAdapter={AdapterDayjs}>
            <RouterProvider router={router} />
          </LocalizationProvider>
        </NotificationsProvider>
      </DialogsProvider>
    </ThemeProvider>
  )
}
