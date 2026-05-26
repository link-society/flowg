import { useColorMode } from '@/theme'

import { useEffect, useState } from 'react'

import Editor, { useMonaco } from '@monaco-editor/react'

import {
  vrlLanguageDefinition,
  vrlThemeDarkDefinition,
  vrlThemeDefinition,
} from '@/lib/vrl-highlighter'

import { VrlCodeEditorProps } from './types'

const VrlCodeEditor = ({ id, code, onCodeChange }: VrlCodeEditorProps) => {
  const [value, setValue] = useState(code)
  const monaco = useMonaco()
  const { mode } = useColorMode()

  useEffect(() => {
    setValue(code)
  }, [code])

  useEffect(() => {
    if (!monaco) return

    monaco.languages.register({ id: 'vrl' })
    monaco.editor.defineTheme('vrl-theme-light', vrlThemeDefinition as any)
    monaco.editor.defineTheme('vrl-theme-dark', vrlThemeDarkDefinition as any)
    monaco.languages.setMonarchTokensProvider(
      'vrl',
      vrlLanguageDefinition as any
    )
    monaco.editor.setTheme(
      mode === 'dark' ? 'vrl-theme-dark' : 'vrl-theme-light'
    )
  }, [monaco, mode])

  const onChange = (val?: string) => {
    setValue(val ?? '')
    onCodeChange(val ?? '')
  }

  return (
    <Editor
      wrapperProps={{ id: id ?? '' }}
      defaultValue={value}
      defaultLanguage="vrl"
      theme={mode === 'dark' ? 'vrl-theme-dark' : 'vrl-theme-light'}
      onChange={onChange}
      options={{ minimap: { enabled: false } }}
    />
  )
}

export default VrlCodeEditor
