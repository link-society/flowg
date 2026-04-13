import { ContainerProps } from '@mui/material'

type AppContainerVariant =
  | 'default'
  | 'page'
  | 'section'
  | 'compact'
  | 'toolbar'

export type AppContainerProps = Omit<ContainerProps, 'maxWidth'> & {
  /**
   * Visual spacing preset
   */
  variant?: AppContainerVariant

  /**
   * Remove horizontal padding
   */
  disableX?: boolean

  /**
   * Remove vertical padding
   */
  disableY?: boolean

  /**
   * Full width container (default true)
   */
  fluid?: boolean
}
