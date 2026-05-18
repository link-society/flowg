import { styled } from '@mui/material'

import AppContainer from '@/components/AppContainer/component'

export const StreamDetailViewContainer = styled(AppContainer)({
  '@media (max-width: 990px)': {
    flexDirection: 'column',
  },
})

export const StreamDetailViewSidebar = styled('div')({
  width: '100%',
  '@media (min-width: 990px)': {
    flex: '0 0 16.67%',
    height: '100%',
  },
})

export const StreamDetailViewContent = styled('div')(({ theme }) => ({
  flex: 1,
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: theme.spacing(1),
  height: '100%',
}))
