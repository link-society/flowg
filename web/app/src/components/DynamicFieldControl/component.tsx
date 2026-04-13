import { useMemo, useState } from 'react'

import InputAdornment from '@mui/material/InputAdornment'
import MenuItem from '@mui/material/MenuItem'
import Select from '@mui/material/Select'
import TextField from '@mui/material/TextField'

import CodeIcon from '@mui/icons-material/Code'
import InputIcon from '@mui/icons-material/Input'

import { DynamicFieldControlProps, EditMode } from './types'

const DynamicFieldControl = <T extends string>(
  props: DynamicFieldControlProps<T>
) => {
  const { value, onChange, select, ...textFieldProps } = props
  const [editMode, setEditMode] = useState<EditMode>(
    value.startsWith('@expr:') ? 'dynamic' : 'static'
  )

  const displayValue = useMemo(() => {
    if (editMode === 'dynamic') {
      return value.slice(6)
    } else {
      return value as T
    }
  }, [editMode, value])

  const updateValue = (newDisplayValue: string) => {
    if (editMode === 'dynamic') {
      onChange(`@expr:${newDisplayValue}`)
    } else {
      onChange(newDisplayValue as T)
    }
  }

  const switchEditMode = (mode: EditMode) => {
    setEditMode((prevMode) => {
      if (prevMode !== mode) {
        if (mode === 'static') {
          onChange('' as T)
        } else if (mode === 'dynamic') {
          onChange(`@expr:`)
        }
      }

      return mode
    })
  }

  return (
    <TextField
      {...textFieldProps}
      select={select && editMode === 'static'}
      value={displayValue}
      onChange={(e) => {
        updateValue(e.target.value)
      }}
      slotProps={{
        input: {
          startAdornment: (
            <InputAdornment position="start">
              <Select<EditMode>
                value={editMode}
                onChange={(e) => {
                  switchEditMode(e.target.value as EditMode)
                }}
                size="small"
              >
                <MenuItem value="static">
                  <InputIcon fontSize="small" />
                </MenuItem>
                <MenuItem value="dynamic">
                  <CodeIcon fontSize="small" />
                </MenuItem>
              </Select>
            </InputAdornment>
          ),
        },
      }}
    />
  )
}

export default DynamicFieldControl
