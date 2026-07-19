import { useTranslation } from 'react-i18next'

import TextField from '@mui/material/TextField'

import { useProfile } from '@/lib/hooks/profile'

const UsernameDisplay = () => {
  const { t } = useTranslation()
  const { user } = useProfile()

  return (
    <TextField
      label={t('components.usernameDisplay.label')}
      value={user.name}
      type="text"
      variant="outlined"
      slotProps={{
        input: {
          readOnly: true,
        },
      }}
    />
  )
}

export default UsernameDisplay
