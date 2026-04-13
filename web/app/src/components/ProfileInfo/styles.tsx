import { Card, CardContent, CardHeader, styled } from '@mui/material'

export const ProfileInfoCard = styled(Card)`
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: stretch;
`

export const ProfileInfoCardHeader = styled(CardHeader)`
  background-color: ${({ theme }) => theme.tokens.colors.headerCardBkg};
  color: ${({ theme }) => theme.tokens.colors.primaryContrast};
  box-shadow: ${({ theme }) => theme.shadows[4]};
`

export const ProfileInfoCardHeaderTitle = styled('div')`
  display: flex;
  align-items: center;
  gap: 0.75rem;
`

export const ProfileInfoCardContent = styled(CardContent)`
  flex: 1 1 0;
  overflow: auto;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  align-items: stretch;
`
