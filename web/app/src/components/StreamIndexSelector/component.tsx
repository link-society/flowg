import Accordion from '@mui/material/Accordion'
import AccordionDetails from '@mui/material/AccordionDetails'
import AccordionSummary from '@mui/material/AccordionSummary'
import Typography from '@mui/material/Typography'

import ExpandMoreIcon from '@mui/icons-material/ExpandMore'

import {
  StreamIndexSelectorChip,
  StreamIndexSelectorContainer,
  StreamIndexSelectorValueList,
} from './styles'
import { StreamIndexSelectorProps, ValueChipProps } from './types'

const ValueChip = (props: ValueChipProps) => {
  const { value, selected, onToggle } = props

  return (
    <StreamIndexSelectorChip
      $selected={selected}
      onClick={() => onToggle(!selected)}
    >
      <Typography variant="caption">{value}</Typography>
    </StreamIndexSelectorChip>
  )
}

const StreamIndexSelector = (props: StreamIndexSelectorProps) => {
  const { indices, selection, onSelectionChange } = props

  const toggleIndexValue = (
    field: string,
    value: string,
    selected: boolean
  ) => {
    const newSelection: Record<string, Array<string>> = {}

    for (const [f, vals] of Object.entries(selection)) {
      newSelection[f] = [...vals]
    }

    if (selected) {
      newSelection[field] = newSelection[field] ?? []
      if (!newSelection[field].includes(value)) {
        newSelection[field].push(value)
      }
    } else {
      newSelection[field] = (newSelection[field] ?? []).filter(
        (v) => v !== value
      )
      if (newSelection[field].length === 0) {
        delete newSelection[field]
      }
    }

    onSelectionChange(newSelection)
  }

  return (
    <StreamIndexSelectorContainer>
      {Object.entries(indices).map(([field, values]) => (
        <Accordion key={field} disableGutters>
          <AccordionSummary expandIcon={<ExpandMoreIcon />}>
            <Typography variant="caption">{field}</Typography>
          </AccordionSummary>
          <AccordionDetails>
            <StreamIndexSelectorValueList>
              {values.map((value) => (
                <ValueChip
                  key={value}
                  value={value}
                  selected={(selection[field] ?? []).includes(value)}
                  onToggle={(isSelected) =>
                    toggleIndexValue(field, value, isSelected)
                  }
                />
              ))}
            </StreamIndexSelectorValueList>
          </AccordionDetails>
        </Accordion>
      ))}
    </StreamIndexSelectorContainer>
  )
}

export default StreamIndexSelector
