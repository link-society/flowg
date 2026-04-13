import Box, { BoxProps } from '@mui/material/Box'
import Typography from '@mui/material/Typography'
import { styled } from '@mui/material/styles'

export const ErrorRoot = styled(Box)({
  padding: '0.75rem',
  display: 'flex',
  flexDirection: 'column',
  gap: '0.75rem',
})

export const ErrorHeading = styled(Typography)(({ theme }) => ({
  fontSize: '1.5rem',
  color: theme.palette.error.main,
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
}))

export const ErrorHeadingLabel = styled(Typography)({
  marginLeft: '0.5rem',
})

export const CodeBlock = styled(Box)<BoxProps<'pre'>>(({ theme }) => ({
  padding: '0.5rem',
  backgroundColor: theme.tokens.colors.black,
  color: theme.tokens.colors.mutedText,
  boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
  fontFamily: 'monospace',
  whiteSpace: 'pre',
  overflow: 'auto',
}))
