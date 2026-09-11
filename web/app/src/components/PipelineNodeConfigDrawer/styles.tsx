import Box from '@mui/material/Box'
import { styled } from '@mui/material/styles'

export const DrawerRoot = styled(Box)(({ theme }) => ({
  flex: '0 0 420px',
  maxWidth: '420px',
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  backgroundColor: theme.palette.background.paper,
  borderLeft: `1px solid ${theme.tokens.colors.borderLight}`,
  boxShadow: theme.shadows[8],
  zIndex: 5,
}))

export const DrawerAccent = styled('div')({
  height: '4px',
  flex: 'none',
})

export const DrawerHeader = styled('div')(({ theme }) => ({
  padding: theme.spacing(1.5),
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: theme.spacing(1.25),
  borderBottom: `1px solid ${theme.tokens.colors.borderLight}`,
  flex: 'none',
}))

export const DrawerHeaderIcon = styled('div')(({ theme }) => ({
  width: '30px',
  height: '30px',
  flex: 'none',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  color: theme.tokens.colors.primaryContrast,
}))

export const DrawerHeaderText = styled('div')({
  minWidth: 0,
  display: 'flex',
  flexDirection: 'column',
  lineHeight: 1.2,
})

export const DrawerHeaderKicker = styled('div')(({ theme }) => ({
  fontSize: '0.6875rem',
  fontWeight: 600,
  letterSpacing: '0.06em',
  textTransform: 'uppercase',
  color: theme.tokens.colors.mutedText,
}))

export const DrawerHeaderTitle = styled('div')({
  fontSize: '1rem',
  fontWeight: 600,
  fontFamily: 'monospace',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

export const DrawerBody = styled('div')(({ theme }) => ({
  flex: '1 1 0',
  minHeight: 0,
  overflow: 'auto',
  padding: theme.spacing(2),
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(2),
}))

export const DrawerActions = styled('div')(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  flexWrap: 'wrap',
  gap: theme.spacing(1),
}))

export const DrawerHint = styled('div')(({ theme }) => ({
  fontSize: '0.8125rem',
  color: theme.tokens.colors.mutedText,
  lineHeight: 1.5,
}))

export const DrawerFooter = styled('div')(({ theme }) => ({
  flex: 'none',
  padding: theme.spacing(1.5),
  borderTop: `1px solid ${theme.tokens.colors.borderLight}`,
  backgroundColor: theme.tokens.colors.codeBg,
  fontSize: '0.75rem',
  lineHeight: 1.5,
  color: theme.tokens.colors.labelText,
  display: 'flex',
  gap: theme.spacing(1),
  alignItems: 'flex-start',
}))
