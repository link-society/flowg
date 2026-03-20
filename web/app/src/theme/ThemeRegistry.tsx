import { CacheProvider } from '@emotion/react'

import { type ReactNode, useMemo } from 'react'

import CssBaseline from '@mui/material/CssBaseline'
import GlobalStyles from '@mui/material/GlobalStyles'
import { ThemeProvider } from '@mui/material/styles'

import { createEmotionCache } from './emotionCache'
import globalStyles from './globalStyles'
import theme from './theme'

interface ThemeRegistryProps {
  readonly children: ReactNode
}

export default function ThemeRegistry({ children }: ThemeRegistryProps) {
  const cache = useMemo(() => createEmotionCache(), [])

  return (
    <CacheProvider value={cache}>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <GlobalStyles styles={globalStyles} />
        {children}
      </ThemeProvider>
    </CacheProvider>
  )
}
