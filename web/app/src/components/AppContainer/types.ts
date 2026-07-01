import { ContainerProps } from '@mui/material'

type AppContainerVariant =
  'default' | 'page' | 'section' | 'compact' | 'toolbar'

export type AppContainerProps = Omit<ContainerProps, 'maxWidth'> & {
  variant?: AppContainerVariant
  disableX?: boolean
  disableY?: boolean
  fluid?: boolean
}
