import { useState } from 'react'

import Button from '@mui/material/Button'
import TextField from '@mui/material/TextField'

import AddIcon from '@mui/icons-material/Add'
import DeleteIcon from '@mui/icons-material/Delete'

type GetKey = (item: string, index: number) => string

type ListEditProps = Readonly<{
  id: string
  list: Array<string>
  setList: (list: Array<string>) => void
  getKey?: GetKey
}>

const defaultKey: GetKey = (_, index) => `${index}`

const ListEdit = ({
  id,
  list,
  setList,
  getKey = defaultKey,
}: ListEditProps) => {
  const [newItem, setNewItem] = useState<string>('')

  return (
    <div className="flex flex-col gap-3">
      {list.map((field, index) => (
        <div
          id={id}
          key={getKey(field, index)}
          className="flex flex-row items-stretch gap-3"
        >
          <TextField
            id={`input:${id}.item.${field}.name`}
            variant="outlined"
            type="text"
            value={field}
            className="grow"
            onChange={({ target }) => {
              setList(list.map((e, i) => (index == i ? target.value : e)))
            }}
          />

          <Button
            id={`btn:${id}.item.${field}.delete`}
            variant="contained"
            color="error"
            size="small"
            onClick={() => {
              setList(list.filter((_, i) => i !== index))
            }}
          >
            <DeleteIcon />
          </Button>
        </div>
      ))}

      <div className="flex flex-row items-stretch gap-3">
        <TextField
          id={`input:${id}.new.name`}
          variant="outlined"
          type="text"
          className="grow"
          value={newItem}
          onChange={({ target }) => setNewItem(target.value)}
        />

        <Button
          id={`btn:${id}.new.add`}
          variant="contained"
          color="primary"
          size="small"
          onClick={() => {
            setList([...list, newItem])
            setNewItem('')
          }}
        >
          <AddIcon />
        </Button>
      </div>
    </div>
  )
}

export default ListEdit
