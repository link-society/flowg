export type DialogConfirmPayload = Readonly<{
  title: string
  message: string
  confirmLabel?: string
  // `danger` is for irreversible destructive actions (deleting a resource).
  // `warning` is for lower-stakes but still noteworthy actions (discarding
  // unsaved local edits, nothing is actually destroyed server-side).
  danger?: boolean
  warning?: boolean
}>
