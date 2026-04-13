import { styled } from '@mui/material'

import Paper from '@mui/material/Paper'

export const StreamIndexSelectorContainer = styled(Paper)`
  height: 100%;
  overflow: auto;
`

export const StreamIndexSelectorValueList = styled('div')`
  display: flex;
  flex-direction: column;
  gap: 4px;
`

export const StreamIndexSelectorChip = styled('button')<{
  $selected?: boolean
}>`
  cursor: pointer;
  padding: 0 8px;
  text-align: left;
  transition: all 150ms ease-in-out;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);

  background-color: ${({ $selected, theme }) =>
    $selected
      ? theme.tokens.colors.selectedBg
      : theme.tokens.colors.disabledBg};
  border: 1px solid
    ${({ $selected, theme }) =>
      $selected
        ? theme.tokens.colors.selectedBorder
        : theme.tokens.colors.borderLight};
  font-weight: ${({ $selected }) => ($selected ? 600 : 400)};
`
