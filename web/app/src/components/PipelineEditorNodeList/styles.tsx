import { Paper, styled } from '@mui/material'

export const NodeListRoot = styled(Paper)({
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
})

export const NodeListHeader = styled('div')(({ theme }) => ({
  padding: '0.5rem',
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  backgroundColor: theme.palette.grey[100],
  boxShadow: theme.shadows[4],
}))

export const NodeListTitle = styled('div')({
  flex: 1,
  fontWeight: 600,
})

export const NodeListLoading = styled('div')({
  flex: '1 1 0',
  minHeight: 0,
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
})

export const NodeListItems = styled('div')({
  flex: '1 1 0',
  minHeight: 0,
  overflow: 'auto',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'flex-start',
  gap: '0.5rem',
  padding: '0.5rem',
})
