import { Box, TextField as MuiTextField, styled } from '@mui/material'

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

export const PipelineDetailViewHeaderLeft = styled('div')(({ theme }) => ({
  display: 'flex',
  flex: 1,
  flexDirection: 'row',
  alignItems: 'center',
  gap: theme.spacing(1.5),
}))

export const PipelineDetailViewHeaderRight = styled('div')(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: theme.spacing(1.5),
}))

export const PipelineDetailViewHeaderTest = styled('div')({
  display: 'flex',
  alignItems: 'center',
})

export const PipelineDetailViewHeaderActions = styled('div')(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: theme.spacing(1.5),
}))

export const PipelineDetailViewBody = styled(AppContainer)(({ theme }) => ({
  position: 'relative',
  alignItems: 'stretch',
  gap: theme.spacing(1),
  overflow: 'auto',
  flex: 1,
  flexDirection: 'column',
  padding: '0 !important',
  [theme.breakpoints.up('md')]: {
    flexDirection: 'row',
    overflow: 'hidden',
  },
}))

export const PipelineDetailViewOverlayHost = styled('div')({
  position: 'absolute',
  inset: 0,
  pointerEvents: 'none',
})

export const PipelineDetailViewLeft = styled('div')(({ theme }) => ({
  width: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  [theme.breakpoints.up('md')]: {
    flex: '0 0 240px',
    height: '100%',
    overflowY: 'auto',
  },
}))

export const PipelineDetailViewLeftItem = styled('div')({
  flex: 'none',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
})

export const PipelineDetailViewCenter = styled('div')({
  flex: 1,
  height: '100%',
})

export const TestDialogHint = styled('div')({
  marginBottom: '0.5rem',
})

export const HeaderNameInput = styled(MuiTextField)(({ theme }) => ({
  '& .MuiOutlinedInput-root': {
    color: theme.tokens.colors.primaryContrast,
    backgroundColor: `rgba(0, 0, 0, ${theme.tokens.opacity.disabled})`,
    '& fieldset': {
      borderColor: theme.tokens.colors.toolbarInputBorder,
    },
  },
  '& .MuiFormLabel-root': {
    color: theme.tokens.colors.primaryContrast,
    '&.Mui-focused': {
      color: theme.tokens.colors.primaryContrast,
    },
  },
}))
