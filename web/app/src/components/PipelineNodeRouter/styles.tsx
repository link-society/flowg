import Box from '@mui/material/Box'
import { styled } from '@mui/material/styles'

import { NodeToolbar } from '@xyflow/react'

export const ToolbarRow = styled(NodeToolbar)({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: 8,
})

export const NodeRoot = styled(Box)(({ theme }) => ({
  width: 270,
  height: 100,
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'stretch',
  gap: 8,
  backgroundColor: theme.tokens.colors.white,
  border: `4px solid ${theme.tokens.colors.nodeRouterBorder}`,
  boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)',
  transition: 'box-shadow 150ms ease-in-out',
  '&:hover': {
    boxShadow:
      '0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)',
  },
}))

export const NodeIcon = styled(Box)(({ theme }) => ({
  backgroundColor: theme.tokens.colors.nodeRouterBg,
  color: theme.tokens.colors.white,
  padding: 12,
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
}))

export const NodeBody = styled(Box)({
  padding: 12,
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
})

export const handleStyle = { width: 12, height: 12 }
