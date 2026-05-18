import { styled } from '@mui/material'

import AppContainer from '@/components/AppContainer/component'

export const StreamDetailViewContainer = styled(AppContainer)(({ theme }) => ({
  [theme.breakpoints.down('md')]: {
    flexDirection: 'column',
  },
}))

export const StreamDetailViewSidebar = styled('div')(({ theme }) => ({
  width: '100%',
  [theme.breakpoints.up('md')]: {
    flex: '0 0 16.67%',
    height: '100%',
  },
}))

export const StreamDetailViewContent = styled('div')(({ theme }) => ({
  flex: 1,
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: theme.spacing(1),
  height: '100%',
}))
