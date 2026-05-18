import { styled } from '@mui/material'

import AppContainer from '@/components/AppContainer/component'

export const AccountViewContainer = styled(AppContainer)(({ theme }) => ({
  display: 'flex',
  height: '100%',
  padding: theme.spacing(1),
  gap: theme.spacing(1),
  [theme.breakpoints.down('lg')]: {
    flexDirection: 'column',
    overflow: 'auto',
  },
}))

export const AccountViewPanel = styled('div')(({ theme }) => ({
  flex: 1,
  height: '100%',
  [theme.breakpoints.down('lg')]: {
    height: 'auto',
  },
}))
