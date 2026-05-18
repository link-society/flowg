import { styled } from '@mui/material'

import Paper from '@mui/material/Paper'

export const StreamIndexSelectorContainer = styled(Paper)({
  height: '100%',
  overflow: 'auto',
})

export const StreamIndexSelectorValueList = styled('div')(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(0.5),
}))

export const StreamIndexSelectorChip = styled('button')<{
  $selected?: boolean
}>(({ $selected, theme }) => ({
  cursor: 'pointer',
  padding: theme.spacing(0, 1),
  textAlign: 'left',
  transition: `all ${theme.tokens.transitions.normal}`,
  boxShadow: `0 1px 2px rgba(0, 0, 0, ${theme.tokens.opacity.light})`,
  backgroundColor: $selected
    ? theme.tokens.colors.selectedBg
    : theme.tokens.colors.disabledBg,
  border: `1px solid ${
    $selected
      ? theme.tokens.colors.selectedBorder
      : theme.tokens.colors.borderLight
  }`,
  fontWeight: $selected ? 600 : 400,
}))
