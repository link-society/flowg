import { styled } from '@mui/material'

export const LogQueryPanelContainer = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  gap: theme.spacing(2),
  padding: theme.spacing(1.5),
}))

export const LogQueryPanelFilterForm = styled('form')({
  flex: 3,
})

export const LogQueryPanelTimeWindow = styled('div')({
  flex: 2,
})

export const LogQueryPanelActions = styled('div')({
  flex: 1,
})
