import { Paper, styled } from '@mui/material'

export const FlowRoot = styled(Paper)({
  width: '100%',
  height: '100%',
})

export const FlowPanelPaper = styled(Paper)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: theme.spacing(1.5),
  boxShadow: theme.shadows[1],
}))

export const FlowPanelLabel = styled('div')(({ theme }) => ({
  alignSelf: 'stretch',
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  fontSize: '0.75rem',
  fontWeight: 600,
  backgroundColor: theme.palette.grey[100],
  borderRight: `1px solid ${theme.palette.grey[200]}`,
  padding: theme.spacing(1),
}))

export const FlowPanelChips = styled('div')(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: theme.spacing(1),
  padding: theme.spacing(1),
}))
