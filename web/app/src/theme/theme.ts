import { colors } from '@mui/material'

import { createTheme } from '@mui/material/styles'

import { type TokensType, tokens } from './tokens'

declare module '@mui/material/styles' {
  interface Theme {
    tokens: TokensType
  }
  interface ThemeOptions {
    tokens?: TokensType
  }
}

const theme = createTheme({
  tokens,
  shape: {
    borderRadius: 0,
  },
  palette: {
    primary: {
      main: colors.blue[800],
    },
    secondary: {
      main: colors.teal[400],
    },
  },
})

export default theme
