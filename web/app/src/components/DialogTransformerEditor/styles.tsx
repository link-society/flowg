import AppBar from '@mui/material/AppBar'
import Box from '@mui/material/Box'
import Toolbar from '@mui/material/Toolbar'
import { styled } from '@mui/material/styles'

export const EditorToolbar = styled(Toolbar)(({ theme }) => ({
  gap: theme.spacing(1.5),
}))

export const TitleField = styled(Box)({
  flexGrow: 1,
})

export const FullScreenBody = styled(Box)(({ theme }) => ({
  flexGrow: 1,
  backgroundColor: theme.tokens.colors.bodyBg,
  padding: theme.spacing(1),
}))

export const FallbackContainer = styled(Box)({
  width: '100%',
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
})

export const DialogAppBar = styled(AppBar)({ position: 'relative' })
