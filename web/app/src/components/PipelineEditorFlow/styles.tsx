import { Chip, Paper, styled } from '@mui/material'

export const FlowRoot = styled(Paper)({
  width: '100%',
  height: '100%',
  '.react-flow__panel': {
    backgroundColor: 'transparent',
  },
})

export const FlowPanelPaper = styled(Paper)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: theme.spacing(1.5),
  boxShadow: theme.shadows[1],
}))

export const FlowPanelLabel = styled('div')(({ theme }) => ({
  alignSelf: 'stretch',
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  fontSize: '0.75rem',
  fontWeight: 600,
  backgroundColor: theme.tokens.colors.cardHeaderBkg,
  color: theme.tokens.colors.primaryContrast,
  borderRight: `1px solid ${theme.tokens.colors.borderLight}`,
  padding: theme.spacing(1),
}))

export const FlowPanelChips = styled('div')(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: theme.spacing(1),
  padding: theme.spacing(1),
}))

export const SwitchNodeChip = styled(Chip)(({ theme }) => ({
  backgroundColor: theme.tokens.colors.switchChipBg,
  borderColor: theme.tokens.colors.switchChipBorder,
  borderRadius: 0,
  boxShadow: theme.shadows[1],
  fontFamily: 'monospace',
  '&:hover': { boxShadow: theme.shadows[4] },
}))
