import {
  Box,
  AppBar as MuiAppBar,
  TextField as MuiTextField,
  Toolbar,
  styled,
} from '@mui/material'

export const DialogAppBar = styled(MuiAppBar)({
  position: 'relative',
})

export const DialogToolbar = styled(Toolbar)(({ theme }) => ({
  gap: theme.spacing(1.5),
}))

export const DialogToolbarName = styled('div')({
  flexGrow: 1,
})

export const ToolbarNameInput = styled(MuiTextField)(({ theme }) => ({
  '& .MuiOutlinedInput-root': {
    color: theme.tokens.colors.primaryContrast,
    backgroundColor: `rgba(0, 0, 0, ${theme.tokens.opacity.disabled})`,
    '& fieldset': {
      borderColor: theme.tokens.colors.toolbarInputBorder,
    },
  },
  '& .MuiFormLabel-root': {
    color: theme.tokens.colors.primaryContrast,
    '&.Mui-focused': {
      color: theme.tokens.colors.primaryContrast,
    },
  },
}))

export const DialogBody = styled(Box)(({ theme }) => ({
  flex: 1,
  padding: theme.spacing(3),
  overflow: 'auto',
}))

export const DialogLoading = styled(Box)({
  width: '100%',
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
})
