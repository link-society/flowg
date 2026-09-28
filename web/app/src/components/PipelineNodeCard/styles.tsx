import Box from '@mui/material/Box'
import { styled } from '@mui/material/styles'

import { NodeToolbar } from '@xyflow/react'

export const ToolbarRow = styled(NodeToolbar)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: theme.spacing(1),
}))

export const NodeRoot = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'borderColor',
})<{ borderColor: string }>(({ theme, borderColor }) => ({
  width: 270,
  height: 100,
  position: 'relative',
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'stretch',
  cursor: 'default',
  backgroundColor: theme.palette.background.paper,
  border: `4px solid ${borderColor}`,
  boxShadow: theme.tokens.shadows.nodeElevated,
  transition: theme.tokens.transitions.shadow,
  '&:hover': {
    boxShadow: theme.tokens.shadows.nodeElevatedHover,
  },
  '&:hover .pipeline-node-drag-handle': {
    opacity: 1,
    pointerEvents: 'auto',
  },
  '&::before': {
    content: '""',
    position: 'absolute',
    top: 0,
    right: '100%',
    width: 8,
    height: '100%',
  },
}))

export const NodeIcon = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'bgColor',
})<{ bgColor: string }>(({ theme, bgColor }) => ({
  backgroundColor: bgColor,
  color: theme.tokens.colors.primaryContrast,
  padding: theme.spacing(1.5),
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  marginRight: theme.spacing(1),
}))

export const NodeBody = styled(Box)(({ theme }) => ({
  padding: theme.spacing(1.5),
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
}))

export const handleStyle = { width: 12, height: 12 }
