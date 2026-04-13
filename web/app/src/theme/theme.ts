import { colors as muiColors } from '@mui/material'

import { createTheme } from '@mui/material/styles'

import { type TokensType, tokens } from './tokens'

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

const theme = createTheme({
  tokens,
  shape: {
    borderRadius: 0,
  },
  palette: {
    primary: {
      main: muiColors.blue[800],
    },
    secondary: {
      main: muiColors.teal[400],
    },
  },
  typography: {
    allVariants: { color: tokens.colors.black },
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
    },
    MuiInputLabel: {
      styleOverrides: {
        root: {
          color: 'rgba(0, 0, 0, 0.6)',
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
  },
})

export default theme
