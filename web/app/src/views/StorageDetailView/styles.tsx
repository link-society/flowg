import { Box, styled } from '@mui/material'

import AppContainer from '@/components/AppContainer/component'

export const StorageDetailViewRoot = styled(Box)`
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: stretch;
`

export const StorageDetailViewHeader = styled(AppContainer)`
  display: flex;
  flex-direction: row;
  align-items: stretch;
  background-color: ${({ theme }) => theme.tokens.colors.toolbarBkg};
  color: ${({ theme }) => theme.tokens.colors.primaryContrast};
  box-shadow: ${({ theme }) => theme.shadows[4]};
  z-index: 10;
`

export const StorageDetailViewHeaderLeft = styled('div')`
  display: flex;
  flex: 1;
  flex-direction: row;
  align-items: center;
  gap: 0.75rem;
`

export const StorageDetailViewHeaderActions = styled('div')`
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 0.75rem;
`

export const StorageDetailViewBody = styled(AppContainer)`
  display: flex;
  flex-direction: row;
  align-items: stretch;
  gap: 0.5rem;
  flex: 1 1 0;
  overflow: hidden;
`

export const StorageDetailViewSidebar = styled('div')`
  width: calc(100% / 6);
  height: 100%;
  flex-shrink: 0;
`

export const StorageDetailViewContent = styled('div')`
  flex: 1;
  height: 100%;
  min-width: 0;
`
