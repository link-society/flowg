import { createContext } from 'react'

import { CloseDialog, OpenDialog } from '@/lib/models/Dialog'

const DialogsContext = createContext<{
  open: OpenDialog
  close: CloseDialog
}>(null!)

export default DialogsContext
