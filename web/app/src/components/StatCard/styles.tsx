import { CardContent, Typography, styled } from '@mui/material'

export const StatCardHeaderWrapper = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  gap: theme.spacing(1.5),
  fontSize: '1.5rem',
  lineHeight: '2rem',
  fontWeight: 600,
}))

export const StatCardContent = styled(CardContent)(({ theme }) => ({
  padding: 0,
  textAlign: 'center',
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(1.25),
}))

export const StatCardValue = styled(Typography)({
  fontWeight: 700,
})
