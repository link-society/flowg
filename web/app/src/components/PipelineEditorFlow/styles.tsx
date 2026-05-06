import { Paper, styled } from '@mui/material'

export const FlowRoot = styled(Paper)`
  width: 100%;
  height: 100%;
`

export const FlowPanelPaper = styled(Paper)`
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 0.75rem;
  box-shadow: ${({ theme }) => theme.shadows[1]};
`

export const FlowPanelLabel = styled('div')`
  align-self: stretch;
  display: flex;
  flex-direction: row;
  align-items: center;
  font-size: 0.75rem;
  font-weight: 600;
  background-color: ${({ theme }) => theme.palette.grey[100]};
  border-right: 1px solid ${({ theme }) => theme.palette.grey[200]};
  padding: 0.5rem;
`

export const FlowPanelChips = styled('div')`
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem;
`
