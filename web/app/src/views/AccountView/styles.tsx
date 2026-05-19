import { styled } from '@mui/material'

import AppContainer from '@/components/AppContainer/component'

export const AccountViewContainer = styled(AppContainer)({
  display: 'flex',
  height: '100%',
  padding: '8px',
  gap: '8px',
  '@media (max-width: 990px)': {
    flexDirection: 'column',
    overflow: 'auto',
  },
})

export const AccountViewPanel = styled('div')({
  flex: 1,
  height: '100%',
  '@media (max-width: 990px)': {
    height: 'auto',
  },
})
