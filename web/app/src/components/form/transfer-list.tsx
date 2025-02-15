import { ReactNode, useEffect, useState } from 'react'

import KeyboardArrowLeftIcon from '@mui/icons-material/KeyboardArrowLeft'
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight'
import KeyboardDoubleArrowLeftIcon from '@mui/icons-material/KeyboardDoubleArrowLeft'
import KeyboardDoubleArrowRightIcon from '@mui/icons-material/KeyboardDoubleArrowRight'

import Grid from '@mui/material/Grid2'
import List from '@mui/material/List'
import ListItemButton from '@mui/material/ListItemButton'
import ListItemIcon from '@mui/material/ListItemIcon'
import ListItemText from '@mui/material/ListItemText'
import Checkbox from '@mui/material/Checkbox'
import Button from '@mui/material/Button'
import Paper from '@mui/material/Paper'

function not<T>(a: readonly T[], b: readonly T[], getItemId: (item: T) => any) {
  return a.filter((value) => !b.find(
    (item) => getItemId(item) === getItemId(value)
  ))
}

function intersection<T>(a: readonly T[], b: readonly T[], getItemId: (item: T) => any) {
  return a.filter((value) => b.find(
    (item) => getItemId(item) === getItemId(value)
  ))
}

type TransferListProps<T> = Readonly<{
  choices: T[]
  getItemId: (item: T) => string
  renderItem: (item: T) => ReactNode
  onChoiceUpdate: (choices: readonly T[]) => void
}>

export function TransferList<T>(props: TransferListProps<T>) {
  const [checked, setChecked] = useState<T[]>([])
  const [left, setLeft] = useState<T[]>(props.choices)
  const [right, setRight] = useState<T[]>([])

  useEffect(
    () => props.onChoiceUpdate(right),
    [right]
  )

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
    <Paper variant="outlined" className="min-h-60 max-h-60 overflow-auto">
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
                  checked={checked.find((v) => {
                    return props.getItemId(v) === props.getItemId(value)
                  }) !== undefined}
                  tabIndex={-1}
                  disableRipple
                  inputProps={{
                    'aria-labelledby': labelId,
                  }}
                />
              </ListItemIcon>
              <ListItemText id={labelId} primary={props.renderItem(value)} />
            </ListItemButton>
          )
        })}
      </List>
    </Paper>
  )

  return (
    <Grid
      container
      spacing={2}
      className="justify-center items-center"
    >
      <Grid
        data-ref="container:generic.transfer-list.items-left"
        size="grow"
      >
        {customList(left)}
      </Grid>
      <Grid>
        <Grid container direction="column" sx={{ alignItems: 'center' }}>
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
        </Grid>
      </Grid>
      <Grid
        data-ref="container:generic.transfer-list.items-right"
        size="grow"
      >
        {customList(right)}
      </Grid>
    </Grid>
  )
}
