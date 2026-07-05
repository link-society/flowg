import { CacheProvider } from '@emotion/react'

import {
  type ReactNode,
  createContext,
  useCallback,
  useContext,
  useEffect,
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
    return globalThis.matchMedia('(prefers-color-scheme: dark)').matches
      ? 'dark'
      : 'light'
  })

  const toggle = useCallback(() => {
    setMode((prev) => {
      const next: ColorMode = prev === 'light' ? 'dark' : 'light'
      localStorage.setItem('colorMode', next)
      return next
    })
  }, [])

  useEffect(() => {
    const mediaQuery = globalThis.matchMedia('(prefers-color-scheme: dark)')
    const handleChange = (e: MediaQueryListEvent) => {
      const saved = localStorage.getItem('colorMode')
      if (!saved) {
        setMode(e.matches ? 'dark' : 'light')
      }
    }
    mediaQuery.addEventListener('change', handleChange)
    return () => mediaQuery.removeEventListener('change', handleChange)
  }, [])

  const theme = useMemo(() => createAppTheme(mode), [mode])
  const colorModeValue = useMemo(() => ({ mode, toggle }), [mode, toggle])

  return (
    <ColorModeContext.Provider value={colorModeValue}>
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
