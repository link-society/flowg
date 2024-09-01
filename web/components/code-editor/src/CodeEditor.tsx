import React, { useCallback, useEffect, useState } from 'react'

import Editor, { useMonaco } from '@monaco-editor/react'

interface CodeEditorProps {
  code: string
  onCodeChange: (value: string) => void
}

const CodeEditor: React.FC<CodeEditorProps> = ({ code, onCodeChange }) => {
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
      // todo...
    },
    [monaco],
  )

  const onChange = useCallback(
    (val?: string) => {
      setValue(val ?? '')
      onCodeChange(val ?? '')
    },
    [onCodeChange],
  )

  return (
    <div className="w-full h-full">
      <Editor
        defaultValue={value}
        width='100%'
        height='100%'
        onChange={onChange}
        options={{minimap: {enabled: false}}}
      />
    </div>
  )
}

export default CodeEditor
