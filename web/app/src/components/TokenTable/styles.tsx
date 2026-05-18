import { Card, CardContent, CardHeader, styled } from '@mui/material'

export const TokenTableCard = styled(Card)({
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  '@media (max-width: 1200px)': {
    minHeight: '24rem',
  },
})

export const TokenTableCardHeader = styled(CardHeader)(({ theme }) => ({
  backgroundColor: theme.tokens.colors.cardHeaderBkg,
  color: theme.tokens.colors.primaryContrast,
  boxShadow: theme.shadows[4],
  zIndex: 20,
}))

export const TokenTableCardHeaderTitle = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  gap: theme.spacing(1.5),
}))

export const TokenTableCardHeaderTitleText = styled('span')({
  flex: 1,
})

export const TokenTableCardContent = styled(CardContent)({
  padding: '0 !important',
  flex: '1 1 0',
})

export const TokenCellRoot = styled('span')({
  fontFamily: 'monospace',
})
