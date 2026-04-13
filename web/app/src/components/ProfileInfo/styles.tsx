import { Card, CardContent, CardHeader, styled } from '@mui/material'

export const ProfileInfoCard = styled(Card)({
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
})

export const ProfileInfoCardHeader = styled(CardHeader)(({ theme }) => ({
  backgroundColor: theme.tokens.colors.cardHeaderBkg,
  color: theme.tokens.colors.primaryContrast,
  boxShadow: theme.shadows[4],
}))

export const ProfileInfoCardHeaderTitle = styled('div')({
  display: 'flex',
  alignItems: 'center',
  gap: '0.75rem',
})

export const ProfileInfoCardContent = styled(CardContent)({
  flex: '1 1 0',
  overflow: 'auto',
  display: 'flex',
  flexDirection: 'column',
  gap: '0.75rem',
  alignItems: 'stretch',
})
