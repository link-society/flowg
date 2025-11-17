import { useEffect, useState } from 'react'

import Button from '@mui/material/Button'
import TextField from '@mui/material/TextField'

import AddIcon from '@mui/icons-material/Add'
import DeleteIcon from '@mui/icons-material/Delete'

type KeyValue = [string, string]

type KeyValueEditorProps = Readonly<{
  id?: string
  keyLabel?: string
  valueLabel?: string
  keyValues: KeyValue[]
  onChange: (kvs: KeyValue[]) => void
}>

export const KeyValueEditor = (props: KeyValueEditorProps) => {
  const [pairs, setPairs] = useState(props.keyValues)

  const [newKey, setNewKey] = useState('')
  const [newValue, setNewValue] = useState('')

  useEffect(() => {
    props.onChange(pairs)
  }, [pairs])

  return (
    <div
      id={props.id ?? 'field:generic.kv-editor'}
      className="flex flex-col items-stretch gap-2"
    >
      {pairs.map(([key, value], index) => (
        <div
          data-ref={`entry:generic.kv-editor.item.${key.toLowerCase().replaceAll(/\s+/g, '-')}`}
          key={key}
          className="flex flex-row items-stretch gap-2"
        >
          <TextField
            data-ref="input:generic.kv-editor.item.key"
            label={props.keyLabel ?? 'Key'}
            value={key}
            onChange={(e) => {
              setPairs((prev) => {
                const next = [...prev]
                next[index] = [e.target.value, value]
                return next
              })
            }}
            variant="outlined"
            size="small"
            className="grow"
          />

          <TextField
            data-ref="input:generic.kv-editor.item.value"
            label={props.valueLabel ?? 'Value'}
            value={value}
            onChange={(e) => {
              setPairs((prev) => {
                const next = [...prev]
                next[index] = [key, e.target.value]
                return next
              })
            }}
            variant="outlined"
            size="small"
            className="grow"
          />

          <Button
            data-ref="btn:generic.kv-editor.item.delete"
            color="error"
            variant="contained"
            size="small"
            onClick={() => {
              setPairs((prev) => {
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
          setPairs((prev) => {
            const newEntry: [string, string] = [newKey, newValue]
            const next = [...prev, newEntry]
            setNewKey('')
            setNewValue('')
            return next
          })
        }}
      >
        <TextField
          data-ref="input:generic.kv-editor.new.key"
          label={props.keyLabel ?? 'Key'}
          value={newKey}
          onChange={(e) => {
            setNewKey(e.target.value)
          }}
          variant="outlined"
          size="small"
          className="grow"
          required
        />

        <TextField
          data-ref="input:generic.kv-editor.new.value"
          label={props.valueLabel ?? 'Value'}
          value={newValue}
          onChange={(e) => {
            setNewValue(e.target.value)
          }}
          variant="outlined"
          size="small"
          className="grow"
          required
        />

        <Button
          data-ref="btn:generic.kv-editor.new.submit"
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
