import { Card, CardContent, CardHeader, styled } from '@mui/material'

export const UserTableCard = styled(Card)(({ theme }) => ({
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  [theme.breakpoints.down('lg')]: {
    minHeight: '24rem',
  },
}))

export const UserTableCardHeader = styled(CardHeader)(({ theme }) => ({
  backgroundColor: theme.tokens.colors.cardHeaderBkg,
  color: theme.tokens.colors.primaryContrast,
  boxShadow: theme.shadows[4],
  zIndex: 20,
}))

export const UserTableCardHeaderTitle = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  gap: theme.spacing(1.5),
}))

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

export const RolesCellRoot = styled('div')(({ theme }) => ({
  display: 'flex',
  flexWrap: 'wrap',
  gap: theme.spacing(0.5),
  padding: theme.spacing(1, 0),
  alignContent: 'center',
  width: '100%',
}))
