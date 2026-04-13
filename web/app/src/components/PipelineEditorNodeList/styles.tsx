import { Paper, styled } from '@mui/material'

export const NodeListRoot = styled(Paper)`
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: stretch;
`

export const NodeListHeader = styled('div')`
  padding: 0.5rem;
  display: flex;
  flex-direction: row;
  align-items: center;
  background-color: ${({ theme }) => theme.palette.grey[100]};
  box-shadow: ${({ theme }) => theme.shadows[4]};
`

export const NodeListTitle = styled('div')`
  flex: 1;
  font-weight: 600;
`

export const NodeListLoading = styled('div')`
  flex: 1 1 0;
  min-height: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
`

export const NodeListItems = styled('div')`
  flex: 1 1 0;
  min-height: 0;
  overflow: auto;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.5rem;
  padding: 0.5rem;
`
