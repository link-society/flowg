import { CardContent, styled } from '@mui/material'

export const StatCardHeaderWrapper = styled('div')({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  gap: '0.75rem',
  fontSize: '1.5rem',
  lineHeight: '2rem',
  fontWeight: 600,
})

export const StatCardContent = styled(CardContent)({
  padding: 0,
  textAlign: 'center',
  display: 'flex',
  flexDirection: 'column',
  gap: 10,
})
