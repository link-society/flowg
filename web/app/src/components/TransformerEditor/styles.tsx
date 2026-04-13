import Box from '@mui/material/Box'
import Paper, { PaperProps } from '@mui/material/Paper'
import Typography from '@mui/material/Typography'
import { styled } from '@mui/material/styles'

export const EditorRoot = styled(Box)({
  display: 'flex',
  flexDirection: 'row',
  width: '100%',
  height: '100%',
  gap: '0.5rem',
})

export const EditorMain = styled(Box)({
  flex: '0 0 calc(66.666% - 0.25rem)',
  minWidth: 0,
})

export const EditorSide = styled(Box)({
  flex: '0 0 calc(33.333% - 0.25rem)',
  minWidth: 0,
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: '0.5rem',
})

export const EditorPaper = styled(Paper)({
  width: '100%',
  height: '100%',
})

export const SidePanel = styled(Box)({
  flex: 1,
  height: 0,
})

export const SidePanelInner = styled(Paper)({
  padding: '0.5rem',
  height: '100%',
  overflow: 'auto',
})

export const SidePanelOutput = styled(Paper)({
  padding: '0.5rem',
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
})

export const SideLabel = styled(Typography)(({ theme }) => ({
  fontSize: '0.875rem',
  color: theme.palette.text.secondary,
  fontWeight: 600,
  marginBottom: '0.5rem',
}))

export const RunButtonRow = styled(Box)({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
})

export const ResultBox = styled(Paper)<PaperProps<'pre'>>(({ theme }) => ({
  padding: '0.5rem',
  flex: 1,
  flexShrink: 1,
  height: 0,
  overflow: 'auto',
  fontFamily: 'monospace',
  backgroundColor: theme.tokens.colors.codeBg,
  margin: 0,
}))
