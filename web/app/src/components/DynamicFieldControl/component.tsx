import { useMemo, useState } from 'react'

import Divider from '@mui/material/Divider'
import FormControl from '@mui/material/FormControl'
import FormHelperText from '@mui/material/FormHelperText'
import InputLabel from '@mui/material/InputLabel'
import MenuItem from '@mui/material/MenuItem'

import CodeIcon from '@mui/icons-material/Code'
import InputIcon from '@mui/icons-material/Input'

import {
  DynamicFieldGroupBorder,
  DynamicFieldGroupRoot,
  DynamicFieldInput,
  DynamicFieldModeSelect,
  DynamicFieldValueSelect,
} from './styles'
import { BorderState, DynamicFieldControlProps, EditMode } from './types'

const DynamicFieldControl = <T extends string>(
  props: DynamicFieldControlProps<T>
) => {
  const {
    value,
    onChange,
    select,
    children,
    error = false,
    disabled = false,
    fullWidth,
    label,
    required,
    id,
    helperText,
    multiline,
    rows,
    minRows,
    maxRows,
    placeholder,
    type,
    autoFocus,
    sx,
    className,
  } = props

  const [editMode, setEditMode] = useState<EditMode>(
    value.startsWith('@expr:') ? 'dynamic' : 'static'
  )
  const [focused, setFocused] = useState(false)

  const displayValue = useMemo(
    () => (editMode === 'dynamic' ? value.slice(6) : (value as T)),
    [editMode, value]
  )

  const shrunk = true
  const borderState: BorderState = { focused, error, disabled, shrunk }
  const isSelect = select && editMode === 'static'

  const updateValue = (newVal: string) => {
    onChange(editMode === 'dynamic' ? `@expr:${newVal}` : (newVal as T))
  }

  const switchEditMode = (mode: EditMode) => {
    setEditMode((prev) => {
      if (prev !== mode) {
        onChange(mode === 'dynamic' ? `@expr:` : ('' as T))
      }
      return mode
    })
  }

  return (
    <FormControl
      error={error}
      disabled={disabled}
      fullWidth={fullWidth}
      focused={focused}
      sx={sx}
      className={className}
    >
      <InputLabel htmlFor={id} shrink={shrunk} required={required}>
        {label}
      </InputLabel>

      <DynamicFieldGroupRoot>
        <DynamicFieldGroupBorder ownerState={borderState} aria-hidden="true">
          <legend>
            <span>{label}</span>
          </legend>
        </DynamicFieldGroupBorder>

        <DynamicFieldModeSelect
          value={editMode}
          onChange={(e) => switchEditMode(e.target.value as EditMode)}
          onOpen={() => setFocused(true)}
          onClose={() => setFocused(false)}
          disabled={disabled}
          inputProps={{ 'aria-label': 'field mode' }}
        >
          <MenuItem value="static">
            <InputIcon fontSize="small" />
          </MenuItem>
          <MenuItem value="dynamic">
            <CodeIcon fontSize="small" />
          </MenuItem>
        </DynamicFieldModeSelect>

        <Divider orientation="vertical" flexItem />

        {isSelect ? (
          <DynamicFieldValueSelect
            id={id}
            value={displayValue}
            onChange={(e) => updateValue(e.target.value as string)}
            onOpen={() => setFocused(true)}
            onClose={() => setFocused(false)}
            onFocus={() => setFocused(true)}
            onBlur={() => setFocused(false)}
            disabled={disabled}
          >
            {children}
          </DynamicFieldValueSelect>
        ) : (
          <DynamicFieldInput
            id={id}
            value={displayValue}
            onChange={(e) => updateValue(e.target.value)}
            onFocus={() => setFocused(true)}
            onBlur={() => setFocused(false)}
            disabled={disabled}
            multiline={multiline}
            rows={rows}
            minRows={minRows}
            maxRows={maxRows}
            placeholder={placeholder}
            type={type}
            autoFocus={autoFocus}
          />
        )}
      </DynamicFieldGroupRoot>

      {helperText && <FormHelperText>{helperText}</FormHelperText>}
    </FormControl>
  )
}

export default DynamicFieldControl
