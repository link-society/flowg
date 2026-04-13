import Box from '@mui/material/Box'
import Paper from '@mui/material/Paper'
import Typography from '@mui/material/Typography'
import { styled } from '@mui/material/styles'

export const StreamEditorRoot = styled(Box)({
  height: '100%',
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'stretch',
  gap: '0.75rem',
})

export const StreamEditorPanel = styled(Paper)({
  flex: 1,
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
})

export const StreamEditorPanelHeader = styled(Box)(({ theme }) => ({
  padding: '0.75rem',
  backgroundColor: theme.palette.grey[100],
  boxShadow: theme.shadows[1],
  textAlign: 'center',
}))

export const StreamEditorPanelBody = styled(Box)({
  flex: '1 1 0',
  height: 0,
  overflow: 'auto',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: '0.75rem',
  padding: '0.75rem',
})

export const StreamEditorUsageRow = styled(Box)({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: '2.5rem',
  marginBottom: '0.5rem',
})

export const StreamEditorHint = styled(Typography)({
  fontStyle: 'italic',
})
