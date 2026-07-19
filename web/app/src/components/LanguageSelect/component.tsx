import { useState } from 'react'
import { useTranslation } from 'react-i18next'

import IconButton from '@mui/material/IconButton'
import ListItemText from '@mui/material/ListItemText'
import Menu from '@mui/material/Menu'
import MenuItem from '@mui/material/MenuItem'
import Typography from '@mui/material/Typography'

import Check from '@mui/icons-material/Check'
import LanguageIcon from '@mui/icons-material/Language'

import { AVAILABLE_LANGUAGES, LANGUAGE_STORAGE_KEY } from '@/lib/i18n'

import { LanguageSelectContainer } from './styles'

const LanguageSelect = () => {
  const { t, i18n } = useTranslation()

  const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null)
  const open = Boolean(anchorEl)

  const currentLanguage =
    AVAILABLE_LANGUAGES.find(
      (language) => language.code === i18n.resolvedLanguage
    ) ?? AVAILABLE_LANGUAGES[0]

  const handleSelect = (code: string) => {
    i18n.changeLanguage(code)
    localStorage.setItem(LANGUAGE_STORAGE_KEY, code)
    setAnchorEl(null)
  }

  return (
    <LanguageSelectContainer>
      <IconButton
        id="btn:footer.toggle-language"
        size="small"
        onClick={(e) => setAnchorEl(e.currentTarget)}
        aria-label={t('components.languageSelect.ariaLabel', {
          label: currentLanguage.label,
        })}
        aria-haspopup="menu"
        aria-expanded={open}
      >
        <LanguageIcon fontSize="small" />
      </IconButton>
      <Typography variant="text" sx={{ fontWeight: 600 }}>
        {currentLanguage.label}
      </Typography>

      <Menu
        anchorEl={anchorEl}
        open={open}
        onClose={() => setAnchorEl(null)}
        anchorOrigin={{ vertical: 'top', horizontal: 'left' }}
        transformOrigin={{ vertical: 'bottom', horizontal: 'left' }}
      >
        {AVAILABLE_LANGUAGES.map((language) => (
          <MenuItem
            key={language.code}
            id={`opt:footer.language-${language.code}`}
            selected={language.code === currentLanguage.code}
            onClick={() => handleSelect(language.code)}
          >
            <ListItemText>{language.label}</ListItemText>
            {language.code === currentLanguage.code && (
              <Check fontSize="small" />
            )}
          </MenuItem>
        ))}
      </Menu>
    </LanguageSelectContainer>
  )
}

export default LanguageSelect
