import { useCallback } from 'react'
import { useReactFlow } from '@xyflow/react'

import DeleteIcon from '@mui/icons-material/Delete'
import Button from '@mui/material/Button'

type DeleteNodeButtonProps = {
  nodeId: string
}

export const DeleteNodeButton = ({ nodeId }: DeleteNodeButtonProps) => {
  const flow = useReactFlow()

  const onDelete = useCallback(
    () => {
      const node = flow.getNode(nodeId)
      if (node !== undefined) {
        const edges = flow.getEdges().filter(
          edge => edge.source === nodeId || edge.target === nodeId
        )
        flow.deleteElements({
          nodes: [node],
          edges
        })
      }
    },
    [nodeId],
  )

  return (
    <Button
      variant="contained"
      size="small"
      color="error"
      onClick={onDelete}
      startIcon={<DeleteIcon />}
    >
      Delete
    </Button>
  )
}
