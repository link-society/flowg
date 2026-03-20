import type { GlobalStylesProps } from '@mui/material/GlobalStyles'

import { colors } from './tokens'

const globalStyles: GlobalStylesProps['styles'] = {
  'html, body': {
    margin: 0,
    padding: 0,
    width: '100vw',
    height: '100vh',
    backgroundColor: colors.backgroundBody,
    overflow: 'hidden',
  },
  '#root': {
    width: '100vw',
    height: '100vh',
    display: 'flex',
    flexDirection: 'column',
  },
  a: {
    textDecoration: 'none',
    color: 'inherit',
  },
}

export default globalStyles
