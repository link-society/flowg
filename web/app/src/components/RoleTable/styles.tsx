import { Card, CardContent, CardHeader, styled } from '@mui/material'

export const RoleTableCard = styled(Card)({
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  '@media (max-width: 1200px)': {
    minHeight: '24rem',
  },
})

export const RoleTableCardHeader = styled(CardHeader)(({ theme }) => ({
  backgroundColor: theme.tokens.colors.cardHeaderBkg,
  color: theme.tokens.colors.primaryContrast,
  boxShadow: theme.shadows[4],
  zIndex: 20,
}))

export const RoleTableCardHeaderTitle = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  gap: theme.spacing(1.5),
}))

export const RoleTableCardHeaderTitleText = styled('span')({
  flex: 1,
})

export const RoleTableCardContent = styled(CardContent)({
  padding: '0 !important',
  flex: '1 1 0',
  overflow: 'hidden',
  '@media (max-width: 990px)': {
    overflowX: 'auto',
  },
  '& .ag-cell-wrapper': {
    height: 'auto !important',
  },
  '& .ag-cell': {
    display: 'flex',
    alignItems: 'center',
  },
  '& .flowg-actions-cell': {
    justifyContent: 'center',
  },
})

export const ScopesCellRoot = styled('div')(({ theme }) => ({
  display: 'flex',
  flexWrap: 'wrap',
  gap: theme.spacing(0.5),
  padding: theme.spacing(1, 0),
  alignContent: 'center',
  width: '100%',
}))
