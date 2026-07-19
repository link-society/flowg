import { type ColorPreference, useColorMode } from '@/theme'

import { useState } from 'react'
import { useTranslation } from 'react-i18next'

import IconButton from '@mui/material/IconButton'
import ListItemIcon from '@mui/material/ListItemIcon'
import ListItemText from '@mui/material/ListItemText'
import Menu from '@mui/material/Menu'
import MenuItem from '@mui/material/MenuItem'
import Typography from '@mui/material/Typography'

import Check from '@mui/icons-material/Check'
import DarkModeIcon from '@mui/icons-material/DarkMode'
import LightModeIcon from '@mui/icons-material/LightMode'
import SettingsBrightnessIcon from '@mui/icons-material/SettingsBrightness'

import { ColorSelectContainer } from './styles'

const PREFERENCE_META = {
  light: {
    icon: LightModeIcon,
    labelKey: 'components.pageFooter.lightMode',
  },
  dark: {
    icon: DarkModeIcon,
    labelKey: 'components.pageFooter.darkMode',
  },
  system: {
    icon: SettingsBrightnessIcon,
    labelKey: 'components.pageFooter.systemDefault',
  },
} as const

const PREFERENCE_ORDER: readonly ColorPreference[] = ['light', 'dark', 'system']

const ColorSelect = () => {
  const { t } = useTranslation()
  const { preference, setPreference } = useColorMode()

  const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null)
  const open = Boolean(anchorEl)

  const { icon: PreferenceIcon, labelKey } = PREFERENCE_META[preference]
  const label = t(labelKey)

  const handleSelect = (next: ColorPreference) => {
    setPreference(next)
    setAnchorEl(null)
  }

  return (
    <ColorSelectContainer>
      <IconButton
        id="btn:footer.toggle-color-mode"
        size="small"
        onClick={(e) => setAnchorEl(e.currentTarget)}
        aria-label={t('components.pageFooter.colorModeAriaLabel', { label })}
        aria-haspopup="menu"
        aria-expanded={open}
      >
        <PreferenceIcon fontSize="small" />
      </IconButton>
      <Typography variant="text" sx={{ fontWeight: 600 }}>
        {label}
      </Typography>

      <Menu
        anchorEl={anchorEl}
        open={open}
        onClose={() => setAnchorEl(null)}
        anchorOrigin={{ vertical: 'top', horizontal: 'left' }}
        transformOrigin={{ vertical: 'bottom', horizontal: 'left' }}
      >
        {PREFERENCE_ORDER.map((value) => {
          const { icon: OptionIcon, labelKey: optionLabelKey } =
            PREFERENCE_META[value]
          const optionLabel = t(optionLabelKey)
          return (
            <MenuItem
              key={value}
              id={`opt:footer.color-mode-${value}`}
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
    </ColorSelectContainer>
  )
}

export default ColorSelect
