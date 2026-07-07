import { type ColorPreference, useColorMode } from '@/theme'

import { useState } from 'react'

import IconButton from '@mui/material/IconButton'
import ListItemIcon from '@mui/material/ListItemIcon'
import ListItemText from '@mui/material/ListItemText'
import Menu from '@mui/material/Menu'
import MenuItem from '@mui/material/MenuItem'

import Check from '@mui/icons-material/Check'
import DarkModeIcon from '@mui/icons-material/DarkMode'
import LightModeIcon from '@mui/icons-material/LightMode'
import SettingsBrightnessIcon from '@mui/icons-material/SettingsBrightness'

import { useFeatureFlags } from '@/lib/hooks/featureflags'

import { PageFooterContainer, PageFooterText } from './styles'

const PREFERENCE_META = {
  light: { icon: LightModeIcon, label: 'Light mode' },
  dark: { icon: DarkModeIcon, label: 'Dark mode' },
  system: { icon: SettingsBrightnessIcon, label: 'System default' },
} as const

const PREFERENCE_ORDER: readonly ColorPreference[] = ['light', 'dark', 'system']

const PageFooter = () => {
  const { demoMode } = useFeatureFlags()
  const { preference, setPreference } = useColorMode()

  const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null)
  const open = Boolean(anchorEl)

  const { icon: PreferenceIcon, label } = PREFERENCE_META[preference]

  const handleSelect = (next: ColorPreference) => {
    setPreference(next)
    setAnchorEl(null)
  }

  return (
    <PageFooterContainer>
      <IconButton
        id="btn:footer.toggle-color-mode"
        size="small"
        onClick={(e) => setAnchorEl(e.currentTarget)}
        aria-label={`Color mode: ${label}. Click to change.`}
        aria-haspopup="menu"
        aria-expanded={open}
      >
        <PreferenceIcon fontSize="small" />
      </IconButton>
      <PageFooterText variant="text">{label}</PageFooterText>

      <Menu
        anchorEl={anchorEl}
        open={open}
        onClose={() => setAnchorEl(null)}
        anchorOrigin={{ vertical: 'top', horizontal: 'left' }}
        transformOrigin={{ vertical: 'bottom', horizontal: 'left' }}
      >
        {PREFERENCE_ORDER.map((value) => {
          const { icon: OptionIcon, label: optionLabel } =
            PREFERENCE_META[value]
          return (
            <MenuItem
              key={value}
              selected={value === preference}
              onClick={() => handleSelect(value)}
            >
              <ListItemIcon>
                <OptionIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText>{optionLabel}</ListItemText>
              {value === preference && <Check fontSize="small" />}
            </MenuItem>
          )
        })}
      </Menu>

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
