import {
  vrlLanguageDefinition,
  vrlThemeDefinition,
} from '@/lib/vrl-highlighter'

import { useEffect, useState } from 'react'

import Editor, { useMonaco } from '@monaco-editor/react'

type VrlCodeEditorProps = Readonly<{
  id?: string
  code: string
  onCodeChange: (value: string) => void
}>

const VrlCodeEditor = ({ id, code, onCodeChange }: VrlCodeEditorProps) => {
  const [value, setValue] = useState(code)
  const monaco = useMonaco()

  useEffect(() => {
    setValue(code)
  }, [code])

  useEffect(() => {
    if (!monaco) return

    monaco.languages.register({ id: 'vrl' })
    monaco.editor.defineTheme('vrl-theme', vrlThemeDefinition as any)
    monaco.languages.setMonarchTokensProvider(
      'vrl',
      vrlLanguageDefinition as any
    )
  }, [monaco])

  const onChange = (val?: string) => {
    setValue(val ?? '')
    onCodeChange(val ?? '')
  }

  return (
    <Editor
      wrapperProps={{ id: id ?? '' }}
      defaultValue={value}
      defaultLanguage="vrl"
      theme="vrl-theme"
      onChange={onChange}
      options={{ minimap: { enabled: false } }}
    />
  )
}

export default VrlCodeEditor
