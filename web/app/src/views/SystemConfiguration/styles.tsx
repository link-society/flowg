import Box from '@mui/material/Box'
import Card from '@mui/material/Card'
import { styled } from '@mui/material/styles'

import AppContainer from '@/components/AppContainer/component'

export const SystemConfigurationRoot = styled(AppContainer)({
  flexDirection: 'column',
})

export const SystemConfigurationHeader = styled(Box)({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
})

export const SystemConfigurationCard = styled(Card)({
  display: 'flex',
  flexDirection: 'column',
})

export const SystemConfigurationWrapper = styled('div')({
  display: 'flex',
  flexDirection: 'column',
  maxWidth: 400,
  width: '100%',
  gap: '1rem',
})

export const SystemConfigurationCardHeader = styled(Box)(({ theme }) => ({
  padding: '0.75rem 1rem',
  backgroundColor: theme.tokens.colors.cardHeaderBkg,
  color: theme.tokens.colors.primaryContrast,
  boxShadow: theme.shadows[4],
  zIndex: 20,
  display: 'flex',
  alignItems: 'center',
  gap: '0.75rem',
}))
