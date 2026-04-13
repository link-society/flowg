import { TextFieldProps } from '@mui/material/TextField'

import { DynamicField } from '@/lib/models/DynamicField'

export type DynamicFieldControlProps<T extends string> = Omit<
  TextFieldProps,
  'value' | 'onChange'
> & {
  value: DynamicField<T>
  onChange: (value: DynamicField<T>) => void
}

export type EditMode = 'static' | 'dynamic'
