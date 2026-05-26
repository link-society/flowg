import { CacheProvider } from '@emotion/react'

import {
  type ReactNode,
  createContext,
  useContext,
  useMemo,
  useState,
} from 'react'

import CssBaseline from '@mui/material/CssBaseline'
import GlobalStyles from '@mui/material/GlobalStyles'
import { ThemeProvider } from '@mui/material/styles'

import { createEmotionCache } from './emotionCache'
import globalStyles from './globalStyles'
import { createAppTheme } from './theme'

interface ThemeRegistryProps {
  readonly children: ReactNode
}

type ColorMode = 'light' | 'dark'

interface ColorModeContextValue {
  mode: ColorMode
  toggle: () => void
}

export const ColorModeContext = createContext<ColorModeContextValue>({
  mode: 'light',
  toggle: () => {},
})

export function useColorMode() {
  return useContext(ColorModeContext)
}

export default function ThemeRegistry({ children }: ThemeRegistryProps) {
  const cache = useMemo(() => createEmotionCache(), [])

  const [mode, setMode] = useState<ColorMode>(() => {
    const saved = localStorage.getItem('colorMode')
    if (saved === 'dark' || saved === 'light') return saved
    return window.matchMedia('(prefers-color-scheme: dark)').matches
      ? 'dark'
      : 'light'
  })

  const toggle = () => {
    setMode((prev) => {
      const next: ColorMode = prev === 'light' ? 'dark' : 'light'
      localStorage.setItem('colorMode', next)
      return next
    })
  }

  const theme = useMemo(() => createAppTheme(mode), [mode])

  return (
    <ColorModeContext.Provider value={{ mode, toggle }}>
      <CacheProvider value={cache}>
        <ThemeProvider theme={theme}>
          <CssBaseline />
          <GlobalStyles styles={globalStyles} />
          {children}
        </ThemeProvider>
      </CacheProvider>
    </ColorModeContext.Provider>
  )
}
