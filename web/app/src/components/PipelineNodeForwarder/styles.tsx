import Box from '@mui/material/Box'
import { styled } from '@mui/material/styles'

import { NodeToolbar } from '@xyflow/react'

export const ToolbarRow = styled(NodeToolbar)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: theme.spacing(1),
}))

export const NodeRoot = styled(Box)(({ theme }) => ({
  width: 270,
  height: 100,
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'stretch',
  gap: theme.spacing(1),
  backgroundColor: theme.tokens.colors.white,
  border: `4px solid ${theme.tokens.colors.nodeForwarderBorder}`,
  boxShadow: theme.tokens.shadows.nodeElevated,
  transition: theme.tokens.transitions.shadow,
  '&:hover': {
    boxShadow: theme.tokens.shadows.nodeElevatedHover,
  },
}))

export const NodeIcon = styled(Box)(({ theme }) => ({
  backgroundColor: theme.tokens.colors.nodeForwarderBg,
  color: theme.tokens.colors.primaryContrast,
  padding: theme.spacing(1.5),
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
}))

export const NodeBody = styled(Box)(({ theme }) => ({
  padding: theme.spacing(1.5),
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
}))

export const handleStyle = { width: 12, height: 12 }
