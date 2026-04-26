import { useContext } from 'react'

import DialogsContext from '@/lib/context/dialogs'

export const useDialogs = () => useContext(DialogsContext)
