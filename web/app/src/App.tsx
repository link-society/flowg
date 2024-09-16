import { RouterProvider } from 'react-router-dom'

import { createTheme, ThemeProvider } from '@mui/material/styles'
import * as colors from '@mui/material/colors'

import { NotificationsProvider } from '@toolpad/core'
import { ConfigProvider } from '@/lib/context/config'

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
    <ConfigProvider
      value={{
        notifications: {
          autoHideDuration: 3000,
        },
      }}
    >
      <ThemeProvider theme={theme}>
        <NotificationsProvider>
          <RouterProvider router={router} />
        </NotificationsProvider>
      </ThemeProvider>
    </ConfigProvider>
  )
}
