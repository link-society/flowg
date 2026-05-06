import { CustomCellRendererProps } from 'ag-grid-react'

export type TableActionsCellProps<T> = CustomCellRendererProps<T> & {
  onDelete?: (data: T) => void
}
