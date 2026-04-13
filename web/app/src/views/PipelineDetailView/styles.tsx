import { Box, styled } from '@mui/material'

import AppContainer from '@/components/AppContainer/component'

export const PipelineDetailViewRoot = styled(Box)`
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: stretch;
`

export const PipelineDetailViewHeader = styled(AppContainer)`
  display: flex;
  flex-direction: row;
  align-items: stretch;
  background-color: ${({ theme }) => theme.tokens.colors.toolbarBkg};
  color: ${({ theme }) => theme.tokens.colors.primaryContrast};
  box-shadow: ${({ theme }) => theme.shadows[4]};
  z-index: 10;
`

export const PipelineDetailViewHeaderLeft = styled('div')`
  display: flex;
  flex: 1;
  flex-direction: row;
  align-items: center;
  gap: 0.75rem;
`

export const PipelineDetailViewHeaderTest = styled('div')`
  display: flex;
  align-items: center;
  margin: 0 0.75rem;
`

export const PipelineDetailViewHeaderActions = styled('div')`
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 0.75rem;
`

export const PipelineDetailViewBody = styled(AppContainer)`
  display: flex;
  flex-direction: row;
  align-items: stretch;
  gap: 0.5rem;
  flex: 1 1 0;
  overflow: hidden;
  padding: 0.5rem;
`

export const PipelineDetailViewLeft = styled('div')`
  width: calc(100% / 6);
  height: 100%;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: stretch;
`

export const PipelineDetailViewCenter = styled('div')`
  flex: 1;
  height: 100%;
  min-width: 0;
`

export const PipelineDetailViewRight = styled('div')`
  width: calc(100% / 6);
  height: 100%;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 0.5rem;
`

export const PipelineDetailViewRightItem = styled('div')`
  flex: 1 1 0;
  min-height: 0;
  display: flex;
  flex-direction: column;
  align-items: stretch;
`

export const TestDialogHint = styled('div')`
  margin-bottom: 0.5rem;
`
