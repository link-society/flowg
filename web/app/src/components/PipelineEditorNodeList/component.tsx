import { ReactElement, ReactNode, useEffect, useState } from 'react'

import Chip from '@mui/material/Chip'
import CircularProgress from '@mui/material/CircularProgress'
import * as colors from '@mui/material/colors'

import OpenInNewIcon from '@mui/icons-material/OpenInNew'

import {
  NodeListHeader,
  NodeListItems,
  NodeListLoading,
  NodeListRoot,
  NodeListTitle,
} from './styles'

type PipelineEditorNodeListProps = Readonly<{
  title: ReactNode
  newButton: (createdCb: () => void) => ReactNode
  fetchItems: () => Promise<string[]>
  itemType: string
  itemIcon: ReactElement
  itemColor: keyof typeof colors
  onItemOpen: (item: string) => void
}>

const PipelineEditorNodeList = (props: PipelineEditorNodeListProps) => {
  const [dirty, setDirty] = useState(true)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  const [items, setItems] = useState<string[]>([])

  const refresh = async () => {
    setLoading(true)
    try {
      const items = await props.fetchItems()
      setItems(items)
      setDirty(false)
    } catch (err) {
      setError(err as Error)
    }

    setLoading(false)
  }

  useEffect(() => {
    refresh()
  }, [dirty])

  useEffect(() => {
    if (error !== null) {
      throw error
    }
  }, [error])

  const baseColorMap = colors[props.itemColor]
  const bgColorIndex = 50 as keyof typeof baseColorMap
  const bdColorIndex = 500 as keyof typeof baseColorMap
  const backgroundColor = baseColorMap[bgColorIndex]
  const borderColor = baseColorMap[bdColorIndex]

  return (
    <NodeListRoot>
      <NodeListHeader>
        <NodeListTitle>{props.title}</NodeListTitle>
        {props.newButton(() => setDirty(true))}
      </NodeListHeader>

      {loading ? (
        <NodeListLoading>
          <CircularProgress size={24} />
        </NodeListLoading>
      ) : (
        <NodeListItems>
          {items.map((item) => (
            <Chip
              key={item}
              icon={props.itemIcon}
              label={item}
              onDelete={() => props.onItemOpen(item)}
              deleteIcon={<OpenInNewIcon />}
              variant="outlined"
              sx={{
                backgroundColor,
                borderColor,
                borderRadius: 0,
                boxShadow: 1,
                fontFamily: 'monospace',
                '&:hover': { boxShadow: 4 },
              }}
              draggable
              onDragStart={(evt) => {
                evt.dataTransfer.setData('item-type', props.itemType)
                evt.dataTransfer.setData('item', item)
                evt.dataTransfer.effectAllowed = 'move'
              }}
            />
          ))}
        </NodeListItems>
      )}
    </NodeListRoot>
  )
}

export default PipelineEditorNodeList
