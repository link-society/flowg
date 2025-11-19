import { useCallback, useEffect, useMemo, useState } from 'react'

import Button from '@mui/material/Button'
import TextField from '@mui/material/TextField'

import AddIcon from '@mui/icons-material/Add'
import DeleteIcon from '@mui/icons-material/Delete'

import * as validators from '@/lib/validators'

type InputListProps = {
  id?: string
  itemLabel?: string
  items: string[]
  itemValidators?: Array<validators.Validator<string>>
  onChange: (items: string[]) => void
}

type Row = {
  id: string
  value: string
}

const genId = () => `id_${Date.now()}_${Math.random().toString(36).slice(2)}`

const InputList = (props: InputListProps) => {
  const propRows = useMemo<Row[]>(
    () => props.items.map((item) => ({ id: genId(), value: item })),
    [props.items]
  )

  const [rows, setRows] = useState(propRows)
  const [newItem, setNewItem] = useState('')

  const rowsValidity = useMemo(
    () =>
      rows.map((row) => {
        for (const validator of props.itemValidators ?? []) {
          if (!validator(row.value)) {
            return false
          }
        }

        return true
      }),
    [rows, props.itemValidators]
  )

  const onNewItemSubmit = useCallback(() => {
    setRows((prev) => {
      const next = [...prev, { id: genId(), value: newItem }]
      setNewItem('')
      return next
    })
  }, [newItem])

  useEffect(() => {
    props.onChange(rows.map((row) => row.value))
  }, [rows])

  return (
    <div
      id={props.id ?? 'field:generic.list-editor'}
      className="flex flex-col items-stretch gap-2"
    >
      {rows.map((row, index) => (
        <div
          data-ref={`entry:generic.list-editor.item.${row.id}`}
          key={row.id}
          className="flex flex-row items-stretch gap-2"
        >
          <TextField
            data-ref="input:generic.list-editor.item"
            label={props.itemLabel ?? 'Item'}
            error={!rowsValidity[index]}
            value={row.value}
            onChange={(e) => {
              setRows((prev) => {
                const next = [...prev]
                next[index] = { ...next[index], value: e.target.value }
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
              setRows((prev) => {
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

      <div className="flex flex-row items-stretch gap-2">
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
          onClick={onNewItemSubmit}
        >
          <AddIcon />
        </Button>
      </div>
    </div>
  )
}

export default InputList
