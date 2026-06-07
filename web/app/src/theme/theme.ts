import { colors as muiColors } from '@mui/material'

import { createTheme } from '@mui/material/styles'

import { type TokensType, tokens } from './tokens'
import { darkColors, lightColors } from './tokens/colors'

declare module '@mui/material/styles' {
  interface Theme {
    tokens: TokensType
  }
  interface ThemeOptions {
    tokens?: TokensType
  }
  interface TypographyVariants {
    titleLg: React.CSSProperties
    titleMd: React.CSSProperties
    titleSm: React.CSSProperties
    text: React.CSSProperties
  }
  interface TypographyVariantsOptions {
    titleLg?: React.CSSProperties
    titleMd?: React.CSSProperties
    titleSm?: React.CSSProperties
    text?: React.CSSProperties
  }
}

declare module '@mui/material/Typography' {
  interface TypographyPropsVariantOverrides {
    titleLg: true
    titleMd: true
    titleSm: true
    text: true
  }
}

declare module '@mui/material/MenuItem' {
  interface MenuItemOwnProps {
    variant?: 'navLink'
  }
}

export function createAppTheme(mode: 'light' | 'dark') {
  const colorTokens = mode === 'dark' ? darkColors : lightColors
  const themeTokens: TokensType = { ...tokens, colors: colorTokens }

  return createTheme({
    tokens: themeTokens,
    shape: {
      borderRadius: 0,
    },
    palette: {
      mode,
      primary: {
        main: mode === 'dark' ? muiColors.blue[400] : muiColors.blue[800],
      },
      secondary: {
        main: mode === 'dark' ? '#4db6ac' : '#26a69a',
      },
      ...(mode === 'dark' && {
        background: {
          default: '#1e1e2e',
          paper: '#27324e',
        },
      }),
    },
    typography: {
      allVariants: { color: colorTokens.labelText },
      titleLg: { fontSize: tokens.typography.titleLg, letterSpacing: 'normal' },
      titleMd: { fontSize: tokens.typography.titleMd, letterSpacing: 'normal' },
      titleSm: { fontSize: tokens.typography.titleSm, letterSpacing: 'normal' },
      text: { fontSize: tokens.typography.text, letterSpacing: 'normal' },
    },
    components: {
      MuiButton: {
        styleOverrides: {
          root: {
            textTransform: 'uppercase',
            letterSpacing: 0,
          },
          contained: {},
          outlined: {},
          text: {},
        },
        variants: [
          {
            props: { variant: 'contained', color: 'secondary' },
            style: { color: '#ffffff' },
          },
        ],
      },
      MuiInputLabel: {
        styleOverrides: {
          root: {
            color:
              mode === 'dark'
                ? 'rgba(255, 255, 255, 0.7)'
                : 'rgba(0, 0, 0, 0.6)',
          },
        },
      },
      MuiCard: {
        styleOverrides: {
          root: {
            height: 'fit-content',
          },
        },
      },
      MuiCardHeader: {
        styleOverrides: {
          root: {
            padding: '12px 16px',
            minHeight: '56px',
          },
          content: {
            overflow: 'hidden',
          },
          title: {
            color: 'inherit',
          },
          subheader: {
            color: 'inherit',
          },
        },
      },
      MuiMenuItem: {
        styleOverrides: {
          root: {
            textTransform: 'uppercase',
          },
        },
        variants: [
          {
            props: { variant: 'navLink' },
            style: ({ theme }) => ({
              color: theme.palette.secondary.main,
              gap: 8,
              '& .MuiSvgIcon-root': {
                color: theme.palette.secondary.main,
              },
            }),
          },
        ],
      },
      MuiList: {
        styleOverrides: {
          root: {
            display: 'flex',
            flexDirection: 'column',
            gap: 8,
          },
        },
      },
      MuiAppBar: {
        styleOverrides: {
          root: {
            backgroundColor: colorTokens.navbarBg,
          },
        },
      },
      MuiPaper: {
        styleOverrides: {
          root: {
            ...(mode === 'dark' && {
              backgroundImage: 'none',
            }),
          },
        },
      },
      MuiMenu: {
        styleOverrides: {
          paper: {
            ...(mode === 'dark' && {
              backgroundColor: '#27324e',
              backgroundImage: 'none',
            }),
          },
        },
      },
    },
  })
}

const theme = createAppTheme('light')

export default theme
