import Box from '@mui/material/Box'
import LinearProgress from '@mui/material/LinearProgress'
import Paper from '@mui/material/Paper'
import Typography from '@mui/material/Typography'
import { styled } from '@mui/material/styles'

export const StreamEditorRoot = styled(Box)(({ theme }) => ({
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: theme.spacing(1.5),

  [theme.breakpoints.up('md')]: {
    flexDirection: 'row',
  },
}))

export const StreamEditorPanel = styled(Paper)({
  flex: 1,
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
})

export const StreamEditorPanelHeader = styled(Box)(({ theme }) => ({
  padding: theme.spacing(1.5),
  backgroundColor: theme.palette.grey[100],
  boxShadow: theme.shadows[1],
  textAlign: 'center',
  '& .MuiTypography-root': { fontWeight: 700 },
}))

export const StreamEditorPanelBody = styled(Box)(({ theme }) => ({
  flex: '1 1 0',
  height: 0,
  overflow: 'auto',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: theme.spacing(1.5),
  padding: theme.spacing(1.5),
}))

export const StreamEditorUsageRow = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: theme.spacing(5),
  marginBottom: '0.5rem',
}))

export const StreamEditorHint = styled(Typography)({
  fontStyle: 'italic',
})

export const StreamEditorProgress = styled(LinearProgress)({
  flexGrow: 1,
  height: '20px',
})
