import { useColorMode } from '@/theme'

import IconButton from '@mui/material/IconButton'

import DarkModeIcon from '@mui/icons-material/DarkMode'
import LightModeIcon from '@mui/icons-material/LightMode'

import { useFeatureFlags } from '@/lib/hooks/featureflags'

import { PageFooterContainer, PageFooterText } from './styles'

const PageFooter = () => {
  const { demoMode } = useFeatureFlags()
  const { mode, toggle } = useColorMode()

  return (
    <PageFooterContainer>
      <IconButton
        id="btn:footer.toggle-color-mode"
        size="small"
        onClick={toggle}
        aria-label={
          mode === 'light' ? 'Switch to dark mode' : 'Switch to light mode'
        }
      >
        {mode === 'light' ? (
          <DarkModeIcon fontSize="small" />
        ) : (
          <LightModeIcon fontSize="small" />
        )}
      </IconButton>
      <PageFooterText variant="text">
        {mode === 'light' ? 'Dark mode' : 'Light mode'}
      </PageFooterText>

      {demoMode && (
        <PageFooterText variant="text">Demo Mode Enabled</PageFooterText>
      )}

      <div>
        <PageFooterText variant="text">
          {import.meta.env.FLOWG_VERSION}
        </PageFooterText>
      </div>
    </PageFooterContainer>
  )
}

export default PageFooter
