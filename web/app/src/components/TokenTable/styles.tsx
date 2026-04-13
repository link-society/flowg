import { Card, CardContent, CardHeader, styled } from '@mui/material'

export const TokenTableCard = styled(Card)`
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: stretch;

  @media (max-width: 1200px) {
    min-height: 24rem;
  }
`

export const TokenTableCardHeader = styled(CardHeader)`
  background-color: ${({ theme }) => theme.tokens.colors.headerCardBkg};
  color: ${({ theme }) => theme.tokens.colors.primaryContrast};
  box-shadow: ${({ theme }) => theme.shadows[4]};
  z-index: 20;
`

export const TokenTableCardHeaderTitle = styled('div')`
  display: flex;
  align-items: center;
  gap: 0.75rem;
`

export const TokenTableCardHeaderTitleText = styled('span')`
  flex: 1;
`

export const TokenTableCardContent = styled(CardContent)`
  padding: 0 !important;
  flex: 1 1 0;
`

export const TokenCellRoot = styled('span')`
  font-family: monospace;
`
