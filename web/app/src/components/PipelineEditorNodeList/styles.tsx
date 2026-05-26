import { Chip, Paper, styled } from '@mui/material'

export const NodeListRoot = styled(Paper)({
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
})

export const NodeListHeader = styled('div')(({ theme }) => ({
  padding: theme.spacing(1),
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  backgroundColor: theme.tokens.colors.cardHeaderBkg,
  color: theme.tokens.colors.primaryContrast,
  boxShadow: theme.shadows[4],
}))

export const NodeListTitle = styled('div')({
  flex: 1,
  fontWeight: 600,
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
