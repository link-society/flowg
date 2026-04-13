import { useEffect, useState } from 'react'

import Button from '@mui/material/Button'
import Checkbox from '@mui/material/Checkbox'
import List from '@mui/material/List'
import ListItemButton from '@mui/material/ListItemButton'
import ListItemIcon from '@mui/material/ListItemIcon'
import ListItemText from '@mui/material/ListItemText'

import KeyboardArrowLeftIcon from '@mui/icons-material/KeyboardArrowLeft'
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight'
import KeyboardDoubleArrowLeftIcon from '@mui/icons-material/KeyboardDoubleArrowLeft'
import KeyboardDoubleArrowRightIcon from '@mui/icons-material/KeyboardDoubleArrowRight'

import {
  TransferColumn,
  TransferControls,
  TransferListPaper,
  TransferRoot,
} from './styles'
import { InputTransferListProps } from './types'

const not = <T,>(
  a: readonly T[],
  b: readonly T[],
  getItemId: (item: T) => any
) => {
  return a.filter(
    (value) => !b.some((item) => getItemId(item) === getItemId(value))
  )
}

const intersection = <T,>(
  a: readonly T[],
  b: readonly T[],
  getItemId: (item: T) => any
) => {
  return a.filter((value) =>
    b.some((item) => getItemId(item) === getItemId(value))
  )
}

const InputTransferList = <T,>(props: InputTransferListProps<T>) => {
  const [checked, setChecked] = useState<T[]>([])
  const [left, setLeft] = useState<T[]>(props.choices)
  const [right, setRight] = useState<T[]>([])

  useEffect(() => props.onChoiceUpdate(right), [right])

  const leftChecked = intersection(checked, left, props.getItemId)
  const rightChecked = intersection(checked, right, props.getItemId)

  const handleToggle = (value: T) => () => {
    const currentIndex = checked.findIndex((v) => {
      return props.getItemId(v) === props.getItemId(value)
    })
    const newChecked = [...checked]

    if (currentIndex === -1) {
      newChecked.push(value)
    } else {
      newChecked.splice(currentIndex, 1)
    }

    setChecked(newChecked)
  }

  const handleAllRight = () => {
    setRight(right.concat(left))
    setLeft([])
  }

  const handleCheckedRight = () => {
    setRight(right.concat(leftChecked))
    setLeft(not(left, leftChecked, props.getItemId))
    setChecked(not(checked, leftChecked, props.getItemId))
  }

  const handleCheckedLeft = () => {
    setLeft(left.concat(rightChecked))
    setRight(not(right, rightChecked, props.getItemId))
    setChecked(not(checked, rightChecked, props.getItemId))
  }

  const handleAllLeft = () => {
    setLeft(left.concat(right))
    setRight([])
  }

  const customList = (items: readonly T[]) => (
    <TransferListPaper variant="outlined">
      <List dense component="div">
        {items.map((value: T) => {
          const itemId = props.getItemId(value)
          const inputId = `checkbox:generic.transfer-list.${itemId}`
          const labelId = `label:generic.transfer-list.${itemId}`

          return (
            <ListItemButton key={itemId} onClick={handleToggle(value)}>
              <ListItemIcon>
                <Checkbox
                  id={inputId}
                  checked={checked.some((v) => {
                    return props.getItemId(v) === props.getItemId(value)
                  })}
                  tabIndex={-1}
                  disableRipple
                  slotProps={{
                    input: {
                      'aria-labelledby': labelId,
                    },
                  }}
                />
              </ListItemIcon>
              <ListItemText id={labelId} primary={props.renderItem(value)} />
            </ListItemButton>
          )
        })}
      </List>
    </TransferListPaper>
  )

  return (
    <TransferRoot>
      <TransferColumn data-ref="container:generic.transfer-list.items-left">
        {customList(left)}
      </TransferColumn>
      <TransferControls>
        <Button
          id="btn:generic.transfer-list.all-right"
          sx={{ my: 0.5 }}
          variant="outlined"
          size="small"
          onClick={handleAllRight}
          disabled={left.length === 0}
          aria-label="move all right"
        >
          <KeyboardDoubleArrowRightIcon />
        </Button>
        <Button
          id="btn:generic.transfer-list.selected-right"
          sx={{ my: 0.5 }}
          variant="outlined"
          size="small"
          onClick={handleCheckedRight}
          disabled={leftChecked.length === 0}
          aria-label="move selected right"
        >
          <KeyboardArrowRightIcon />
        </Button>
        <Button
          id="btn:generic.transfer-list.selected-left"
          sx={{ my: 0.5 }}
          variant="outlined"
          size="small"
          onClick={handleCheckedLeft}
          disabled={rightChecked.length === 0}
          aria-label="move selected left"
        >
          <KeyboardArrowLeftIcon />
        </Button>
        <Button
          id="btn:generic.transfer-list.all-left"
          sx={{ my: 0.5 }}
          variant="outlined"
          size="small"
          onClick={handleAllLeft}
          disabled={right.length === 0}
          aria-label="move all left"
        >
          <KeyboardDoubleArrowLeftIcon />
        </Button>
      </TransferControls>
      <TransferColumn data-ref="container:generic.transfer-list.items-right">
        {customList(right)}
      </TransferColumn>
    </TransferRoot>
  )
}

export default InputTransferList
