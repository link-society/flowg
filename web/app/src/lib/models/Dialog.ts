import { ComponentType } from 'react'

export interface DialogProps<P = undefined, R = void> {
  payload: P
  open: boolean
  onClose: (result: R) => Promise<void>
}

export type DialogComponent<P, R> = ComponentType<DialogProps<P, R>>

export interface OpenDialogOptions<R> {
  onClose?: (result: R) => Promise<void>
}

export interface OpenDialog {
  <P extends undefined, R>(
    Component: DialogComponent<P, R>,
    payload?: P,
    options?: OpenDialogOptions<R>
  ): Promise<R>

  <P, R>(
    Component: DialogComponent<P, R>,
    payload: P,
    options?: OpenDialogOptions<R>
  ): Promise<R>
}

export type CloseDialog = <R>(dialog: Promise<R>, result: R) => Promise<R>

export interface DialogStackEntry<P, R> {
  key: string
  open: boolean
  promise: Promise<R>
  Component: DialogComponent<P, R>
  payload: P
  onClose: (result: R) => Promise<void>
  resolve: (result: R) => void
}
