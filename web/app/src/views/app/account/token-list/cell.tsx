import { CustomCellRendererProps } from 'ag-grid-react'

type TokenCellProps = CustomCellRendererProps<string> & {}

export const TokenCell = (props: TokenCellProps) => (
  <span className="font-mono">{props.value}</span>
)
