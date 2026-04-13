import { styled } from '@mui/material'

import AppContainer from '@/components/AppContainer/component'

export const TransformerDetailViewContainer = styled('div')`
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: stretch;
`

export const TransformerDetailViewToolbar = styled(AppContainer)`
  display: flex;
  flex-direction: row;
  align-items: stretch;
  color: white;
  background-color: ${({ theme }) => theme.tokens.colors.editorToolbarBkg};
  z-index: 10;
  box-shadow:
    0 4px 6px -1px rgba(0, 0, 0, 0.1),
    0 2px 4px -1px rgba(0, 0, 0, 0.06);
`

export const TransformerDetailViewToolbarLeft = styled('div')`
  flex-grow: 1;
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 12px;
`

export const TransformerDetailViewToolbarRight = styled('div')`
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 12px;
`

export const TransformerDetailViewBody = styled(AppContainer)`
  flex-grow: 1;
  flex-shrink: 1;
  height: 0;
  display: flex;
  flex-direction: row;
  gap: 8px;
  /* padding: 8px; */
`

export const TransformerDetailViewSidebar = styled('div')`
  flex: 0 0 16.67%;
  height: 100%;
`

export const TransformerDetailViewEditor = styled('div')`
  flex: 1;
  height: 100%;
`
