import { styled } from '@mui/material'

import AppContainer from '@/components/AppContainer/component'

export const StreamDetailViewContainer = styled(AppContainer)`
  display: flex;
  height: 100%;
  padding: 8px;
  gap: 8px;
`

export const StreamDetailViewSidebar = styled('div')`
  flex: 0 0 16.67%;
  height: 100%;
`

export const StreamDetailViewContent = styled('div')`
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 8px;
`
