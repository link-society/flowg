import { RouterProvider } from 'react-router-dom'
import { SnackbarProvider } from 'notistack'
import { createTheme, ThemeProvider } from '@mui/material/styles'
import * as colors from '@mui/material/colors'

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
      <SnackbarProvider
        maxSnack={3}
        autoHideDuration={3000}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'right',
        }}
      >
        <RouterProvider router={router} />
      </SnackbarProvider>
    </ThemeProvider>
  )
}
