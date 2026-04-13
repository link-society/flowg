import { useState } from 'react'

import Button from '@mui/material/Button'
import TextField from '@mui/material/TextField'

import AddIcon from '@mui/icons-material/Add'
import DeleteIcon from '@mui/icons-material/Delete'

import { Root, Row } from './styles'
import { ListEditProps } from './types'

const defaultKey = (_: string, index: number) => `${index}`

const ListEdit = ({
  id,
  list,
  setList,
  getKey = defaultKey,
}: ListEditProps) => {
  const [newItem, setNewItem] = useState<string>('')

  return (
    <Root>
      {list.map((field, index) => (
        <Row id={id} key={getKey(field, index)}>
          <TextField
            id={`input:${id}.item.${field}.name`}
            variant="outlined"
            type="text"
            value={field}
            sx={{ flexGrow: 1 }}
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
        </Row>
      ))}

      <Row>
        <TextField
          id={`input:${id}.new.name`}
          variant="outlined"
          type="text"
          sx={{ flexGrow: 1 }}
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
      </Row>
    </Root>
  )
}

export default ListEdit
