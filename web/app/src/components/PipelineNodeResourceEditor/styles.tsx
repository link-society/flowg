import Paper from '@mui/material/Paper'
import { styled } from '@mui/material/styles'

export const EditorFrame = styled('div')(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: theme.spacing(1.5),
}))

export const EditorLoading = styled('div')(({ theme }) => ({
  display: 'flex',
  justifyContent: 'center',
  padding: theme.spacing(3),
}))

export const CodeSurface = styled(Paper)({
  height: '280px',
  overflow: 'hidden',
})
