import { styled } from '@mui/material'

import AppContainer from '@/components/AppContainer/component'

export const AdminViewContainer = styled(AppContainer)({
  display: 'flex',
  height: '100%',
  padding: '8px',
  gap: '8px',
  '@media (max-width: 1200px)': {
    flexDirection: 'column',
    overflow: 'auto',
  },
})

export const AdminViewPanel = styled('div')({
  flex: 1,
  height: '100%',
  '@media (max-width: 1200px)': {
    height: 'auto',
  },
  '@media (max-width: 990px)': {
    width: '100%',
  },
})
