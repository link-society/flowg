import Box from '@mui/material/Box'
import { styled } from '@mui/material/styles'

export const IndicatorDot = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'hasError',
})<{ hasError: boolean }>(({ theme, hasError }) => ({
  width: '18px',
  height: '18px',
  position: 'absolute',
  right: '-9px',
  top: '-9px',
  backgroundColor: hasError
    ? theme.tokens.colors.statusError
    : theme.tokens.colors.statusSuccess,
  borderRadius: '50%',
  boxShadow: `-2px 2px 2px ${theme.tokens.colors.shadowDark}`,
}))
