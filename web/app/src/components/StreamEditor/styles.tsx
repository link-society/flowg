import Box from '@mui/material/Box'
import LinearProgress from '@mui/material/LinearProgress'
import Paper from '@mui/material/Paper'
import Typography from '@mui/material/Typography'
import { styled } from '@mui/material/styles'

type StackedProps = { stacked?: boolean }

const withoutStacked = {
  shouldForwardProp: (prop: string) => prop !== 'stacked',
}

export const StreamEditorRoot = styled(
  Box,
  withoutStacked
)<StackedProps>(({ theme, stacked }) => ({
  height: stacked ? 'auto' : '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: theme.spacing(1.5),

  ...(stacked
    ? {}
    : {
        [theme.breakpoints.up('md')]: {
          flexDirection: 'row',
        },
      }),
}))

export const StreamEditorPanel = styled(
  Paper,
  withoutStacked
)<StackedProps>(({ stacked }) => ({
  flex: 1,
  height: stacked ? 'auto' : '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
}))

export const StreamEditorPanelHeader = styled(
  Box,
  withoutStacked
)<StackedProps>(({ theme, stacked }) => ({
  padding: theme.spacing(0.75, 1.5),
  backgroundColor: stacked
    ? theme.tokens.colors.codeBg
    : theme.tokens.colors.cardHeaderBkg,
  color: stacked
    ? theme.tokens.colors.labelText
    : theme.tokens.colors.primaryContrast,
  boxShadow: theme.shadows[1],
  textAlign: 'center',
  '& .MuiTypography-root': { fontWeight: 700 },
}))

export const StreamEditorPanelBody = styled(
  Box,
  withoutStacked
)<StackedProps>(({ theme, stacked }) => ({
  flex: stacked ? 'none' : '1 1 0',
  height: stacked ? 'auto' : 0,
  overflow: stacked ? 'visible' : 'auto',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: theme.spacing(1.5),
  padding: theme.spacing(1.5),
}))

export const StreamEditorUsageRow = styled(
  Box,
  withoutStacked
)<StackedProps>(({ theme, stacked }) => ({
  display: 'flex',
  flexDirection: stacked ? 'column' : 'row',
  alignItems: stacked ? 'stretch' : 'center',
  gap: theme.spacing(stacked ? 1 : 5),
  marginBottom: '0.5rem',
}))

export const StreamEditorHint = styled(Typography)({
  fontStyle: 'italic',
})

export const StreamEditorProgress = styled(LinearProgress)({
  flexGrow: 1,
  height: '20px',
})
