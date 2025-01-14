import { RouterProvider } from 'react-router'

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
      <RouterProvider router={router} />
    </ThemeProvider>
  )
}
