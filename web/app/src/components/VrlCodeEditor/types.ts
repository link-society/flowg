export type VrlCodeEditorProps = Readonly<{
  id?: string
  code: string
  onCodeChange: (value: string) => void
}>
