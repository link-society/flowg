import { useEffect, useState } from 'react'

import Button from '@mui/material/Button'
import TextField from '@mui/material/TextField'

import AddIcon from '@mui/icons-material/Add'
import DeleteIcon from '@mui/icons-material/Delete'

type ListEditorProps = {
  id?: string
  itemLabel?: string
  items: string[]
  onChange: (items: string[]) => void
}


export const ListEditor = (props: ListEditorProps) => {
  const [items, setItems] = useState(props.items)

  const [newItem, setNewItem] = useState('')

  useEffect(() => {
    props.onChange(items)
  }, [items])

  return (
    <div
      id={props.id ?? 'field:generic.list-editor'}
      className="flex flex-col items-stretch gap-2"
    >
      {items.map((item, index) => (
        <div
          data-ref={`entry:generic.list-editor.item.${index}`}
          key={index}
          className="flex flex-row items-stretch gap-2"
        >
          <TextField
            data-ref="input:generic.list-editor.item"
            label={props.itemLabel ?? 'Item'}
            value={item}
            onChange={(e) => {
              setItems((prev) => {
                const next = [...prev]
                next[index] = e.target.value
                return next
              })
            }}
            variant="outlined"
            size="small"
            className="grow"
          />

          <Button
            data-ref="btn:generic.list-editor.item.delete"
            color="error"
            variant="contained"
            size="small"
            onClick={() => {
              setItems((prev) => {
                const next = [...prev]
                next.splice(index, 1)
                return next
              })
            }}
          >
            <DeleteIcon />
          </Button>
        </div>
      ))}

      <form
        className="flex flex-row items-stretch gap-2"
        onSubmit={(e) => {
          e.preventDefault()
          setItems((prev) => {
            const next = [...prev, newItem]
            setNewItem('')
            return next
          })
        }}
      >
        <TextField
          data-ref="input:generic.list-editor.new"
          label={props.itemLabel ?? 'Item'}
          value={newItem}
          onChange={(e) => {
            setNewItem(e.target.value)
          }}
          variant="outlined"
          size="small"
          className="grow"
          required
        />

        <Button
          data-ref="btn:generic.list-editor.new.submit"
          color="primary"
          variant="contained"
          size="small"
          type="submit"
        >
          <AddIcon />
        </Button>
      </form>
    </div>
  )
}
