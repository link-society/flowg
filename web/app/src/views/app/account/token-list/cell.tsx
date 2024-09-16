import { CustomCellRendererProps } from 'ag-grid-react'

import { RowType } from './types'

type TokenCellProps = CustomCellRendererProps<RowType> & {}

export const TokenCell = (props: TokenCellProps) => (
  <span className="font-mono">{props.value}</span>
)
