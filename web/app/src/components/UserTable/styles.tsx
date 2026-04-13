import { Card, CardContent, CardHeader, styled } from '@mui/material'

export const UserTableCard = styled(Card)({
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  '@media (max-width: 1200px)': {
    minHeight: '24rem',
  },
})

export const UserTableCardHeader = styled(CardHeader)(({ theme }) => ({
  backgroundColor: theme.tokens.colors.cardHeaderBkg,
  color: theme.tokens.colors.primaryContrast,
  boxShadow: theme.shadows[4],
  zIndex: 20,
}))

export const UserTableCardHeaderTitle = styled('div')({
  display: 'flex',
  alignItems: 'center',
  gap: '0.75rem',
})

export const UserTableCardHeaderTitleText = styled('span')({
  flex: 1,
})

export const UserTableCardContent = styled(CardContent)({
  padding: '0 !important',
  flex: '1 1 0',
  overflow: 'hidden',
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

export const RolesCellRoot = styled('div')({
  display: 'flex',
  flexWrap: 'wrap',
  gap: 4,
  padding: '8px 0',
  alignContent: 'center',
  width: '100%',
})
