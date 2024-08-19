import React, { useCallback, useEffect, useState } from 'react'

import CodeMirror, { ViewUpdate } from '@uiw/react-codemirror'
import { vscodeLight } from '@uiw/codemirror-theme-vscode'

interface CodeEditorProps {
  code: string
  onCodeChange: (value: string) => void
}

const CodeEditor: React.FC<CodeEditorProps> = ({ code, onCodeChange }) => {
  const [value, setValue] = useState(code)

  useEffect(() => {
    setValue(code)
  }, [code])

  const onChange = useCallback((val: string, viewUpdate: ViewUpdate) => {
    setValue(val)
    onCodeChange(val)
  }, [onCodeChange])

  return (
    <div style={{ width: '100%', height: '100%' }}>
      <CodeMirror
        value={value}
        width='100%'
        height='100%'
        theme={vscodeLight}
        onChange={onChange}
        style={{ height: '100%', overflow: 'auto' }}
      />
    </div>
  )
}

export default CodeEditor
