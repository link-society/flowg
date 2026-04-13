import { styled } from '@mui/material'

import Paper from '@mui/material/Paper'

export const StreamIndexSelectorContainer = styled(Paper)({
  height: '100%',
  overflow: 'auto',
})

export const StreamIndexSelectorValueList = styled('div')({
  display: 'flex',
  flexDirection: 'column',
  gap: 4,
})

export const StreamIndexSelectorChip = styled('button')<{
  $selected?: boolean
}>(({ $selected, theme }) => ({
  cursor: 'pointer',
  padding: '0 8px',
  textAlign: 'left',
  transition: 'all 150ms ease-in-out',
  boxShadow: '0 1px 2px rgba(0, 0, 0, 0.1)',
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
