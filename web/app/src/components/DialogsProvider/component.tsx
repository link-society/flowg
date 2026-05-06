import { useId, useMemo, useRef, useState } from 'react'

import useEventCallback from '@mui/utils/useEventCallback'

import DialogsContext from '@/lib/context/dialogs'

import {
  DialogComponent,
  DialogStackEntry,
  OpenDialog,
  OpenDialogOptions,
} from '@/lib/models/Dialog'

import { DialogsProviderProps } from './types'

const DialogsProvider = (props: DialogsProviderProps) => {
  const { children, unmountAfter = 1000 } = props
  const [stack, setStack] = useState<DialogStackEntry<any, any>[]>([])
  const keyPrefix = useId()
  const nextId = useRef(0)
  const dialogMetadata = useRef(
    new WeakMap<Promise<any>, DialogStackEntry<any, any>>()
  )

  const requestDialog = useEventCallback<OpenDialog>(
    <P, R>(
      Component: DialogComponent<P, R>,
      payload: P,
      options: OpenDialogOptions<R> = {}
    ) => {
      const { onClose = async () => {} } = options
      let resolve: ((result: R) => void) | undefined
      const promise = new Promise<R>((res) => {
        resolve = res
      })
      if (resolve === undefined) {
        throw new Error('resolve not set')
      }

      const key = `${keyPrefix}-${nextId.current}`
      nextId.current += 1

      const newEntry: DialogStackEntry<P, R> = {
        key,
        open: true,
        promise,
        Component,
        payload,
        onClose,
        resolve,
      }

      dialogMetadata.current.set(promise, newEntry)

      setStack((prev) => [...prev, newEntry])

      return promise
    }
  )

  const closeDialogUi = useEventCallback(<R,>(dialog: Promise<R>) => {
    setStack((prev) =>
      prev.map((entry) =>
        entry.promise === dialog ? { ...entry, open: false } : entry
      )
    )

    setTimeout(() => {
      setStack((prev) => prev.filter((entry) => entry.promise !== dialog))
    }, unmountAfter)
  })

  const closeDialog = useEventCallback(
    async <R,>(dialog: Promise<R>, result: R) => {
      const entryToClose = dialogMetadata.current.get(dialog)
      if (!entryToClose) {
        throw new Error('Dialog not found')
      }

      try {
        await entryToClose.onClose(result)
      } finally {
        entryToClose.resolve(result)
        closeDialogUi(dialog)
      }

      return dialog
    }
  )

  const contextValue = useMemo(
    () => ({ open: requestDialog, close: closeDialog }),
    [requestDialog, closeDialog]
  )

  return (
    <DialogsContext.Provider value={contextValue}>
      {children}
      {stack.map(({ key, open, Component, payload, promise }) => (
        <Component
          key={key}
          open={open}
          payload={payload}
          onClose={async (result) => {
            await closeDialog(promise, result)
          }}
        />
      ))}
    </DialogsContext.Provider>
  )
}

export default DialogsProvider
