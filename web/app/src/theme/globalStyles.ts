import type { GlobalStylesProps } from '@mui/material/GlobalStyles'

const globalStyles: GlobalStylesProps['styles'] = (theme) => ({
  'html, body': {
    margin: 0,
    padding: 0,
    width: '100vw',
    height: '100vh',
    backgroundColor: theme.tokens.colors.bodyBg,
    overflow: 'hidden',
  },
  '#root': {
    width: '100vw',
    height: '100vh',
    display: 'flex',
    flexDirection: 'column',
  },
  a: {
    textDecoration: 'none',
    color: 'inherit',
  },
  '.flowg-table .ag-header': {
    backgroundColor: theme.tokens.colors.disabledBg,
    zIndex: 10,
    boxShadow: theme.shadows[1],
  },
  '.flowg-table .ag-cell-wrapper': {
    height: '100%',
  },
  '.flowg-table .flowg-actions-header .ag-header-cell-label': {
    justifyContent: 'center',
  },
})

export default globalStyles
