import { Box, styled } from '@mui/material'

import AppContainer from '@/components/AppContainer/component'

export const PipelineDetailViewRoot = styled(Box)({
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
})

export const PipelineDetailViewHeader = styled(AppContainer)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'stretch',
  backgroundColor: theme.tokens.colors.toolbarBkg,
  color: theme.tokens.colors.primaryContrast,
  boxShadow: theme.shadows[4],
  zIndex: 10,
  flex: 0,
}))

export const PipelineDetailViewHeaderLeft = styled('div')({
  display: 'flex',
  flex: 1,
  flexDirection: 'row',
  alignItems: 'center',
  gap: '0.75rem',
})

export const PipelineDetailViewHeaderTest = styled('div')({
  display: 'flex',
  alignItems: 'center',
  margin: '0 0.75rem',
})

export const PipelineDetailViewHeaderActions = styled('div')({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: '0.75rem',
})

export const PipelineDetailViewBody = styled(AppContainer)({
  alignItems: 'stretch',
  gap: '0.5rem',
  overflow: 'auto',
  flex: 1,
  flexDirection: 'column',
  '@media (min-width: 990px)': {
    flexDirection: 'row',
    overflow: 'hidden',
  },
})

export const PipelineDetailViewLeft = styled('div')({
  flex: '0 0 16.67%',
  height: '100%',
  width: '100%',
  '@media (min-width: 990px)': {
    flexDirection: 'row',
  },
})

export const PipelineDetailViewCenter = styled('div')({
  flex: 1,
  height: '100%',
  minHeight: 600,
})

export const PipelineDetailViewRight = styled('div')({
  width: '100%',
  height: '100%',
  flexShrink: 0,
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: '0.5rem',
  '@media (min-width: 990px)': {
    width: 'calc(100% / 6)',
  },
})

export const PipelineDetailViewRightItem = styled('div')({
  flex: '1 1 0',
  minHeight: 0,
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
})

export const TestDialogHint = styled('div')({
  marginBottom: '0.5rem',
})
