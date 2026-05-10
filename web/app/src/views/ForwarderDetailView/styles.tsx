import { Box, Paper, styled } from '@mui/material'

import AppContainer from '@/components/AppContainer/component'

export const ForwarderDetailViewRoot = styled(Box)({
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
})

export const ForwarderDetailViewHeader = styled(AppContainer)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'stretch',
  backgroundColor: theme.tokens.colors.toolbarBkg,
  color: theme.tokens.colors.primaryContrast,
  boxShadow: theme.shadows[4],
  zIndex: 10,
  flex: 0,
}))

export const ForwarderDetailViewHeaderLeft = styled('div')({
  display: 'flex',
  flex: 1,
  flexDirection: 'row',
  alignItems: 'center',
  gap: '0.75rem',
})

export const ForwarderDetailViewHeaderRight = styled('div')({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: '0.75rem',
})

export const ForwarderDetailViewHeaderTest = styled('div')({
  display: 'flex',
  alignItems: 'center',
})

export const ForwarderDetailViewHeaderActions = styled('div')({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: '0.75rem',
})

export const ForwarderDetailViewBody = styled(AppContainer)({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: '0.5rem',
  flex: 1,
  overflow: 'hidden',
  '@media (min-width: 990px)': {
    flexDirection: 'row',
  },
})

export const ForwarderDetailViewSidebar = styled('div')({
  flex: '0 0 16.67%',
  height: '100%',
  width: '100%',
  '@media (min-width: 990px)': {
    flexDirection: 'row',
  },
})

export const ForwarderDetailViewContent = styled('div')({
  flex: 1,
  height: '100%',
  minWidth: 0,
})

export const ForwarderDetailViewEditorPaper = styled(Paper)({
  height: '100%',
  overflow: 'auto',
  padding: '0.75rem',
})

export const TestDialogHint = styled('div')({
  marginBottom: '0.5rem',
})
