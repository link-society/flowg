import { styled } from '@mui/material'

import AppContainer from '@/components/AppContainer/component'

export const AccountViewContainer = styled(AppContainer)(({ theme }) => ({
  display: 'flex',
  height: '100%',
  padding: theme.spacing(1),
  gap: theme.spacing(1),
  '@media (max-width: 1200px)': {
    flexDirection: 'column',
    overflow: 'auto',
  },
}))

export const AccountViewPanel = styled('div')({
  flex: 1,
  height: '100%',
  '@media (max-width: 1200px)': {
    height: 'auto',
  },
})
