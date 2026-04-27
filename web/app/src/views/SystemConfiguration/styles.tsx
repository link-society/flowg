import Box from '@mui/material/Box'
import Card from '@mui/material/Card'
import { styled } from '@mui/material/styles'

import AppContainer from '@/components/AppContainer/component'

export const SystemConfigurationRoot = styled(AppContainer)({
  width: 'calc(100% / 3)',
  margin: 'auto',
  display: 'flex',
  flexDirection: 'column',
  gap: '0.5rem',
})

export const SystemConfigurationHeader = styled(Box)({
  marginBottom: '1.5rem',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
})

export const SystemConfigurationCard = styled(Card)({
  display: 'flex',
  flexDirection: 'column',
})

export const SystemConfigurationCardHeader = styled(Box)(({ theme }) => ({
  padding: '0.75rem 1rem',
  backgroundColor: theme.tokens.colors.headerCardBkg,
  color: theme.tokens.colors.primaryContrast,
  boxShadow: theme.shadows[4],
  zIndex: 20,
  display: 'flex',
  alignItems: 'center',
  gap: '0.75rem',
}))
