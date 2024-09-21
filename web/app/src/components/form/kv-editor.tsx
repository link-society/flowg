import { useEffect, useState } from 'react'

import DeleteIcon from '@mui/icons-material/Delete'
import AddIcon from '@mui/icons-material/Add'

import TextField from '@mui/material/TextField'
import Button from '@mui/material/Button'

type KeyValue = [string, string]

type KeyValueEditorProps = Readonly<{
  keyLabel?: string
  valueLabel?: string
  keyValues: KeyValue[]
  onChange: (kvs: KeyValue[]) => void
}>

export const KeyValueEditor = (props: KeyValueEditorProps) => {
  const [pairs, setPairs] = useState(props.keyValues)

  const [newKey, setNewKey] = useState('')
  const [newValue, setNewValue] = useState('')

  useEffect(
    () => { props.onChange(pairs) },
    [pairs]
  )

  return (
    <div className="flex flex-col items-stretch gap-2">
      {pairs.map(([key, value], index) => (
        <div
          key={key}
          className="flex flex-row items-stretch gap-2"
        >
          <TextField
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
            className="flex-grow"
          />

          <TextField
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
            className="flex-grow"
          />

          <Button
            color="primary"
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
          label={props.keyLabel ?? 'Key'}
          value={newKey}
          onChange={(e) => {
            setNewKey(e.target.value)
          }}
          variant="outlined"
          size="small"
          className="flex-grow"
          required
        />

        <TextField
          label={props.valueLabel ?? 'Value'}
          value={newValue}
          onChange={(e) => {
            setNewValue(e.target.value)
          }}
          variant="outlined"
          size="small"
          className="flex-grow"
          required
        />

        <Button
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
