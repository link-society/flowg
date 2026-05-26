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
  '.ag-theme-balham': {
    ...(theme.palette.mode === 'dark' && {
      '--ag-background-color': '#27324e',
      '--ag-header-background-color': '#1e284a',
      '--ag-header-cell-hover-background-color': '#263356',
      '--ag-odd-row-background-color': '#243050',
      '--ag-row-hover-color': 'rgba(66, 165, 245, 0.08)',
      '--ag-selected-row-background-color': 'rgba(66, 165, 245, 0.2)',
      '--ag-foreground-color': '#cbd5e1',
      '--ag-header-foreground-color': '#cbd5e1',
      '--ag-secondary-foreground-color': '#94a3b8',
      '--ag-border-color': '#2d3a52',
      '--ag-row-border-color': '#2d3a52',
      '--ag-range-selection-border-color': '#42a5f5',
      '--ag-cell-horizontal-border': 'none',
    }),
  },
  '.ag-theme-material': {
    ...(theme.palette.mode === 'dark' && {
      '--ag-background-color': '#27324e',
      '--ag-header-background-color': '#1e284a',
      '--ag-header-cell-hover-background-color': '#263356',
      '--ag-odd-row-background-color': '#243050',
      '--ag-row-hover-color': 'rgba(66, 165, 245, 0.08)',
      '--ag-selected-row-background-color': 'rgba(66, 165, 245, 0.2)',
      '--ag-foreground-color': '#cbd5e1',
      '--ag-header-foreground-color': '#cbd5e1',
      '--ag-secondary-foreground-color': '#94a3b8',
      '--ag-border-color': '#2d3a52',
      '--ag-row-border-color': '#2d3a52',
      '--ag-range-selection-border-color': '#42a5f5',
      '--ag-cell-horizontal-border': 'none',
    }),
  },
  ...(theme.palette.mode === 'dark' && {
    '.react-flow__controls': {
      backgroundColor: '#27324e',
      border: '1px solid #2d3a52',
      boxShadow: 'none',
    },
    '.react-flow__controls-button': {
      backgroundColor: '#27324e',
      borderBottom: '1px solid #2d3a52',
      color: '#cbd5e1',
      fill: '#cbd5e1',
      '&:hover': {
        backgroundColor: '#1e3a8a',
      },
    },
    [`.ag-header-group-cell:not(.ag-column-resizing) + .ag-header-group-cell:not(.ag-column-hover):not(.ag-header-cell-moving):hover,
      .ag-header-group-cell:not(.ag-column-resizing) + .ag-header-group-cell:not(.ag-column-hover).ag-column-resizing,
      .ag-header-cell:not(.ag-column-resizing) + .ag-header-cell:not(.ag-column-hover):not(.ag-header-cell-moving):hover,
      .ag-header-cell:not(.ag-column-resizing) + .ag-header-cell:not(.ag-column-hover).ag-column-resizing,
      .ag-header-group-cell:first-of-type:not(.ag-header-cell-moving):hover,
      .ag-header-group-cell:first-of-type.ag-column-resizing,
      .ag-header-cell:not(.ag-column-hover):first-of-type:not(.ag-header-cell-moving):hover,
      .ag-header-cell:not(.ag-column-hover):first-of-type.ag-column-resizing`]:
      {
        backgroundColor: 'var(--ag-header-cell-hover-background-color)',
      },
  }),
})

export default globalStyles
