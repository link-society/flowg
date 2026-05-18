import { Box, styled } from '@mui/material'

import AppContainer from '@/components/AppContainer/component'

export const StorageDetailViewRoot = styled(Box)({
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
})

export const StorageDetailViewHeader = styled(AppContainer)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'stretch',
  backgroundColor: theme.tokens.colors.toolbarBkg,
  color: theme.tokens.colors.primaryContrast,
  boxShadow: theme.shadows[4],
  zIndex: 10,
  flex: 0,
}))

export const StorageDetailViewHeaderLeft = styled('div')(({ theme }) => ({
  display: 'flex',
  flex: 1,
  flexDirection: 'row',
  alignItems: 'center',
  gap: theme.spacing(1.5),
}))

export const StorageDetailViewHeaderActions = styled('div')(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: theme.spacing(1.5),
}))

export const StorageDetailViewBody = styled(AppContainer)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: theme.spacing(1),
  flex: '1 1 0',
  overflow: 'hidden',
  '@media (min-width: 990px)': {
    flexDirection: 'row',
  },
}))

export const StorageDetailViewSidebar = styled('div')({
  flex: '0 0 16.67%',
  height: '100%',
  width: '100%',
  '@media (min-width: 990px)': {
    flexDirection: 'row',
  },
})

export const StorageDetailViewContent = styled('div')({
  flex: 1,
  height: '100%',
  minWidth: 0,
})
