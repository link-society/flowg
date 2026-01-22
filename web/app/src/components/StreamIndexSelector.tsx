import Accordion from '@mui/material/Accordion'
import AccordionDetails from '@mui/material/AccordionDetails'
import AccordionSummary from '@mui/material/AccordionSummary'
import Paper from '@mui/material/Paper'
import Typography from '@mui/material/Typography'

import ExpandMoreIcon from '@mui/icons-material/ExpandMore'

type ValueChipProps = {
  value: string
  selected: boolean
  onToggle: (selected: boolean) => void
}

const ValueChip = (props: ValueChipProps) => {
  const { value, selected, onToggle } = props

  if (selected) {
    return (
      <div
        className="
          cursor-pointer
          px-2 py-0
          bg-blue-200
          border border-blue-300
          font-semibold
          shadow-xs
          transition-all duration-150 ease-in-out
        "
        onClick={() => onToggle(false)}
        role="button"
      >
        <Typography variant="caption">{value}</Typography>
      </div>
    )
  } else {
    return (
      <div
        className="
          cursor-pointer
          px-2 py-0
          bg-gray-200
          border border-gray-300
          shadow-xs
          transition-all duration-150 ease-in-out
        "
        onClick={() => onToggle(true)}
        role="button"
      >
        <Typography variant="caption">{value}</Typography>
      </div>
    )
  }
}

type StreamIndexSelectorProps = {
  indices: Record<string, Array<string>>
  selection: Record<string, Array<string>>
  onSelectionChange: (selection: Record<string, Array<string>>) => void
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
    <Paper className="h-full overflow-auto">
      {Object.entries(indices).map(([field, values]) => (
        <Accordion key={field} className="w-full" disableGutters>
          <AccordionSummary expandIcon={<ExpandMoreIcon />}>
            <Typography variant="caption">{field}</Typography>
          </AccordionSummary>
          <AccordionDetails>
            <div className="flex flex-col gap-1">
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
            </div>
          </AccordionDetails>
        </Accordion>
      ))}
    </Paper>
  )
}

export default StreamIndexSelector
