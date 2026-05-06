import Box from '@mui/material/Box'
import Toolbar from '@mui/material/Toolbar'
import { styled } from '@mui/material/styles'

export const EditorToolbar = styled(Toolbar)({
  gap: '0.75rem',
})

export const TitleField = styled(Box)({
  flexGrow: 1,
})

export const FullScreenBody = styled(Box)(({ theme }) => ({
  flexGrow: 1,
  backgroundColor: theme.tokens.colors.backgroundBody,
  padding: '1.5rem',
  overflow: 'auto',
}))

export const FallbackContainer = styled(Box)({
  width: '100%',
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
})
