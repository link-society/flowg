import { Box, Paper, styled } from '@mui/material'

import AppContainer from '@/components/AppContainer/component'

export const ForwarderDetailViewRoot = styled(Box)`
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: stretch;
`

export const ForwarderDetailViewHeader = styled(AppContainer)`
  display: flex;
  flex-direction: row;
  align-items: stretch;
  background-color: ${({ theme }) => theme.tokens.colors.toolbarBkg};
  color: ${({ theme }) => theme.tokens.colors.primaryContrast};
  box-shadow: ${({ theme }) => theme.shadows[4]};
  z-index: 10;
`

export const ForwarderDetailViewHeaderLeft = styled('div')`
  display: flex;
  flex: 1;
  flex-direction: row;
  align-items: center;
  gap: 0.75rem;
`

export const ForwarderDetailViewHeaderTest = styled('div')`
  display: flex;
  align-items: center;
  margin: 0 0.75rem;
`

export const ForwarderDetailViewHeaderActions = styled('div')`
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 0.75rem;
`

export const ForwarderDetailViewBody = styled(AppContainer)`
  display: flex;
  flex-direction: row;
  align-items: stretch;
  gap: 0.5rem;
  flex: 1 1 0;
  overflow: hidden;
`

export const ForwarderDetailViewSidebar = styled('div')`
  width: 16.6667%;
  height: 100%;
  flex-shrink: 0;
`

export const ForwarderDetailViewContent = styled('div')`
  flex: 1;
  height: 100%;
  min-width: 0;
`

export const ForwarderDetailViewEditorPaper = styled(Paper)`
  height: 100%;
  overflow: auto;
  padding: 0.75rem;
`

export const TestDialogHint = styled('div')`
  margin-bottom: 0.5rem;
`
