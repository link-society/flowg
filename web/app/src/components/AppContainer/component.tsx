import { Container } from '@mui/material'

import { forwardRef } from 'react'

import { AppContainerProps } from './types'

export const AppContainer = forwardRef<HTMLDivElement, AppContainerProps>(
  (
    {
      children,
      variant = 'default',
      disableX = false,
      disableY = false,
      fluid = true,
      sx,
      ...props
    },
    ref
  ) => {
    const getPadding = () => {
      switch (variant) {
        case 'section':
          return {
            px: disableX ? 0 : { xs: 2, md: 3 },
            py: 0,
          }
        case 'toolbar':
          return {
            px: disableX ? 0 : { xs: 1, md: 3 },
            py: disableX ? 0 : { xs: 1, md: 1 },
          }
        case 'compact':
          return { px: 0, py: 0 }
        case 'page':
        default:
          return {
            px: disableX ? 0 : { xs: 2, md: 3 },
            py: disableY ? 0 : { xs: 2, md: 3 },
          }
      }
    }

    return (
      <Container
        ref={ref}
        maxWidth={fluid ? false : 'xl'}
        {...props}
        sx={(theme) => ({
          ...getPadding(),
          ...sx,
          display: 'flex',
          placeItems: 'center',
          gap: theme.spacing(4),
          flex: 1,
        })}
      >
        {children}
      </Container>
    )
  }
)

export default AppContainer
