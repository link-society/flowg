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
export type ColorPreference = 'light' | 'dark' | 'system'

interface ColorModeContextValue {
  mode: ColorMode
  preference: ColorPreference
  setPreference: (preference: ColorPreference) => void
}

export const ColorModeContext = createContext<ColorModeContextValue>({
  mode: 'light',
  preference: 'system',
  setPreference: () => {},
})

export function useColorMode() {
  return useContext(ColorModeContext)
}

export default function ThemeRegistry({ children }: ThemeRegistryProps) {
  const cache = useMemo(() => createEmotionCache(), [])

  const [preference, setPreference] = useState<ColorPreference>(() => {
    const saved = localStorage.getItem('colorMode')
    if (saved === 'dark' || saved === 'light' || saved === 'system') return saved
    return 'system'
  })

  const [systemMode, setSystemMode] = useState<ColorMode>(() =>
    globalThis.matchMedia('(prefers-color-scheme: dark)').matches
      ? 'dark'
      : 'light',
  )

  const mode: ColorMode = preference === 'system' ? systemMode : preference

  const setColorPreference = useCallback((next: ColorPreference) => {
    localStorage.setItem('colorMode', next)
    setPreference(next)
  }, [])

  useEffect(() => {
    const mediaQuery = globalThis.matchMedia('(prefers-color-scheme: dark)')
    const handleChange = (e: MediaQueryListEvent) => {
      setSystemMode(e.matches ? 'dark' : 'light')
    }
    mediaQuery.addEventListener('change', handleChange)
    return () => mediaQuery.removeEventListener('change', handleChange)
  }, [])

  const theme = useMemo(() => createAppTheme(mode), [mode])
  const colorModeValue = useMemo(
    () => ({ mode, preference, setPreference: setColorPreference }),
    [mode, preference, setColorPreference],
  )

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
