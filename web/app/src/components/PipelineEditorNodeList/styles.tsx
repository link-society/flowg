import { Chip, Paper, styled } from '@mui/material'

export const NodeListRoot = styled(Paper)({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
})

export const NodeListHeader = styled('div')(({ theme }) => ({
  minHeight: 38,
  padding: theme.spacing(0.5, 1),
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  backgroundColor: theme.tokens.colors.codeBg,
  color: theme.tokens.colors.labelText,
  boxShadow: theme.shadows[4],
}))

export const NodeListHeaderToggle = styled('div')({
  flex: 1,
  minWidth: 0,
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: 4,
  cursor: 'pointer',
})

export const NodeListExpandIcon = styled('div', {
  shouldForwardProp: (prop) => prop !== 'expanded',
})<{ expanded: boolean }>(({ expanded }) => ({
  display: 'flex',
  transform: expanded ? 'rotate(0deg)' : 'rotate(-90deg)',
  transition: 'transform 0.15s ease',
}))

export const NodeListTitle = styled('div')({
  flex: 1,
  minWidth: 0,
  fontSize: '0.8125rem',
  fontWeight: 600,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

export const NodeListLoading = styled('div')({
  flex: '1 1 0',
  minHeight: 0,
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
})

export const NodeListItems = styled('div')(({ theme }) => ({
  flex: '1 1 0',
  minHeight: 0,
  overflow: 'auto',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'flex-start',
  gap: theme.spacing(1),
  padding: theme.spacing(1),
}))

export const NodeChip = styled(Chip, {
  shouldForwardProp: (prop) =>
    prop !== 'chipBgColor' && prop !== 'chipBorderColor',
})<{ chipBgColor: string; chipBorderColor: string }>(
  ({ theme, chipBgColor, chipBorderColor }) => ({
    backgroundColor: chipBgColor,
    borderColor: chipBorderColor,
    borderRadius: 0,
    boxShadow: theme.shadows[1],
    fontFamily: 'monospace',
    '&:hover': { boxShadow: theme.shadows[4] },
  })
)
