import { ReactNode, useEffect, useState } from 'react'

import Chip from '@mui/material/Chip'
import CircularProgress from '@mui/material/CircularProgress'
import Paper from '@mui/material/Paper'
import * as colors from '@mui/material/colors'

import OpenInNewIcon from '@mui/icons-material/OpenInNew'

type NodeListProps = Readonly<{
  title: ReactNode
  newButton: (createdCb: () => void) => ReactNode
  fetchItems: () => Promise<string[]>
  itemType: string
  itemIcon: ReactNode
  itemColor: keyof typeof colors
  onItemOpen: (item: string) => void
  className?: string
}>

export function NodeList(props: NodeListProps) {
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
    <Paper className={props.className}>
      <div className="h-full flex flex-col items-stretch">
        <div className="p-2 flex flex-row items-center bg-gray-100 shadow-md">
          <div className="grow text-semibold">{props.title}</div>
          {props.newButton(() => setDirty(true))}
        </div>
        {loading ? (
          <div
            className="
                grow shrink h-0
                flex flex-col items-center justify-center
              "
          >
            <CircularProgress size={24} />
          </div>
        ) : (
          <div
            className="
                grow shrink h-0 overflow-auto
                flex flex-col items-start gap-2 p-2
              "
          >
            {items.map((item) => (
              <Chip
                key={item}
                icon={<>{props.itemIcon}</>}
                label={item}
                onDelete={() => props.onItemOpen(item)}
                deleteIcon={<OpenInNewIcon />}
                variant="outlined"
                sx={{
                  backgroundColor,
                  borderColor,
                }}
                className="rounded-none! shadow-xs hover:shadow-lg font-mono!"
                draggable
                onDragStart={(evt) => {
                  evt.dataTransfer.setData('item-type', props.itemType)
                  evt.dataTransfer.setData('item', item)
                  evt.dataTransfer.effectAllowed = 'move'
                }}
              />
            ))}
          </div>
        )}
      </div>
    </Paper>
  )
}
