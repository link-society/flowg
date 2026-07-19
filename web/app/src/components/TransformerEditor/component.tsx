import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'

import PlayArrowIcon from '@mui/icons-material/PlayArrow'

import * as testApi from '@/lib/api/operations/tests'

import { useApiOperation } from '@/lib/hooks/api'

import InputKeyValue from '@/components/InputKeyValue/component'
import VrlCodeEditor from '@/components/VrlCodeEditor/component'

import {
  EditorMain,
  EditorPaper,
  EditorRoot,
  EditorSide,
  ResultBox,
  RunButtonRow,
  SideLabel,
  SidePanel,
  SidePanelInner,
  SidePanelOutput,
} from './styles'
import { TransformerEditorProps } from './types'

const TransformerEditor = (props: TransformerEditorProps) => {
  const { t } = useTranslation()
  const [code, setCode] = useState(props.code)
  const [testRecord, setTestRecord] = useState<[string, string][]>([])
  const [testResult, setTestResult] = useState<string>('')

  useEffect(() => {
    props.onCodeChange(code)
  }, [code])

  const [onTest, testLoading] = useApiOperation(async () => {
    const input = Object.fromEntries(testRecord)
    const output = await testApi.testTransformer(code, input)

    if (output.success) {
      setTestResult(JSON.stringify(output.records, null, 2))
    } else {
      setTestResult(output.error)
    }
  }, [code, testRecord])

  return (
    <EditorRoot>
      <EditorMain>
        <EditorPaper>
          <VrlCodeEditor
            id="monaco:transformers.editor"
            code={code}
            onCodeChange={setCode}
          />
        </EditorPaper>
      </EditorMain>
      <EditorSide>
        <SidePanel>
          <SidePanelInner>
            <SideLabel variant="text">
              {t('common.testDialog.inputRecord')}
            </SideLabel>
            <InputKeyValue
              id="kv:transformers.test.record"
              keyLabel={t('common.testDialog.fieldLabel')}
              keyValues={testRecord}
              onChange={setTestRecord}
            />
          </SidePanelInner>
        </SidePanel>
        <RunButtonRow>
          <Button
            id="btn:transformers.test.run"
            variant="contained"
            color="primary"
            endIcon={!testLoading && <PlayArrowIcon />}
            disabled={testLoading}
            onClick={() => onTest()}
          >
            {testLoading ? (
              <CircularProgress size={24} />
            ) : (
              <>{t('common.actions.run')}</>
            )}
          </Button>
        </RunButtonRow>
        <SidePanel>
          <SidePanelOutput>
            <SideLabel variant="text">
              {t('components.transformerEditor.outputRecordLabel')}
            </SideLabel>
            <ResultBox
              id="container:transformers.test.result"
              variant="outlined"
              component="pre"
            >
              {testResult}
            </ResultBox>
          </SidePanelOutput>
        </SidePanel>
      </EditorSide>
    </EditorRoot>
  )
}

export default TransformerEditor
