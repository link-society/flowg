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

export const ProfileInfoCardHeaderTitle = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  gap: theme.spacing(1.5),
}))

export const ProfileInfoCardContent = styled(CardContent)(({ theme }) => ({
  flex: '1 1 0',
  overflow: 'auto',
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(1.5),
  alignItems: 'stretch',
}))
