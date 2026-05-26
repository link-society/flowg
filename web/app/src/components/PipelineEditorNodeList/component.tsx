import { useColorMode } from '@/theme'

import { ReactElement, ReactNode, useEffect, useState } from 'react'

import CircularProgress from '@mui/material/CircularProgress'
import * as colors from '@mui/material/colors'

import OpenInNewIcon from '@mui/icons-material/OpenInNew'

import {
  NodeChip,
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
      setError(err instanceof Error ? err : new Error(String(err)))
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

  const { mode } = useColorMode()

  const baseColorMap = colors[props.itemColor]
  const bgColorIndex = (mode === 'dark' ? 900 : 50) as keyof typeof baseColorMap
  const bdColorIndex = (
    mode === 'dark' ? 300 : 500
  ) as keyof typeof baseColorMap
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
            <NodeChip
              key={item}
              icon={props.itemIcon}
              label={item}
              onDelete={() => props.onItemOpen(item)}
              deleteIcon={<OpenInNewIcon />}
              variant="outlined"
              chipBgColor={backgroundColor}
              chipBorderColor={borderColor}
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
