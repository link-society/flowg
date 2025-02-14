import { useEffect, useState } from 'react'

import Editor, { useMonaco } from '@monaco-editor/react'
import { vrlLanguageDefinition, vrlThemeDefinition } from '@/lib/vrl-highlighter'

type CodeEditorProps = Readonly<{
  code: string
  onCodeChange: (value: string) => void
}>

export const CodeEditor = ({ code, onCodeChange }: CodeEditorProps) => {
  const [value, setValue] = useState(code)
  const monaco = useMonaco()

  useEffect(
    () => {
      setValue(code)
    },
    [code],
  )

  useEffect(
    () => {
      if (!monaco) return

      monaco.languages.register({id: 'vrl'})
      monaco.editor.defineTheme('vrl-theme', vrlThemeDefinition as any)
      monaco.languages.setMonarchTokensProvider('vrl', vrlLanguageDefinition as any)
    },
    [monaco],
  )

  const onChange = (val?: string) => {
    setValue(val ?? '')
    onCodeChange(val ?? '')
  }

  return (
    <Editor
      defaultValue={value}
      defaultLanguage='vrl'
      theme='vrl-theme'
      onChange={onChange}
      options={{minimap: {enabled: false}}}
    />
  )
}
