import { Typography, styled } from '@mui/material'

import SearchOffIcon from '@mui/icons-material/SearchOff'

export const NotFoundViewContainer = styled('div')(({ theme }) => ({
  flex: 1,
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  gap: theme.spacing(2),
  padding: theme.spacing(4),
  textAlign: 'center',
}))

export const NotFoundIcon = styled(SearchOffIcon)({
  fontSize: '5rem',
  opacity: 0.3,
})

export const NotFoundTitle = styled(Typography)({
  fontWeight: 700,
})

export const NotFoundHint = styled(Typography)(({ theme }) => ({
  color: `rgba(0, 0, 0, ${theme.tokens.opacity.overlay})`,
}))
