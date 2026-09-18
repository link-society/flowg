import { ReactNode, useEffect, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'

import CircularProgress from '@mui/material/CircularProgress'

import * as configApi from '@/lib/api/operations/config'
import * as logApi from '@/lib/api/operations/logs'

import { useApiOperation } from '@/lib/hooks/api'
import { useDirty } from '@/lib/hooks/dirty'
import { useNotify } from '@/lib/hooks/notify'

import ForwarderModel from '@/lib/models/ForwarderModel'
import StreamConfigModel from '@/lib/models/StreamConfigModel'

import ForwarderEditor from '@/components/ForwarderEditor/component'
import StreamEditor from '@/components/StreamEditor/component'
import VrlCodeEditor from '@/components/VrlCodeEditor/component'

import { CodeSurface, EditorFrame, EditorLoading } from './styles'
import {
  PipelineNodeResourceEditorProps,
  PipelineNodeResourceSaveAction,
} from './types'

type ResourceFormProps = Readonly<{
  saving: boolean
  canSave: boolean
  onSave: () => void
  onSaveActionChange: (action: PipelineNodeResourceSaveAction | null) => void
  children: ReactNode
}>

const Loading = () => (
  <EditorLoading>
    <CircularProgress size={28} />
  </EditorLoading>
)

const ResourceForm = ({
  saving,
  canSave,
  onSave,
  onSaveActionChange,
  children,
}: ResourceFormProps) => {
  useEffect(() => {
    onSaveActionChange({ saving, canSave, onSave })
    return () => onSaveActionChange(null)
  }, [saving, canSave, onSave, onSaveActionChange])

  return <EditorFrame>{children}</EditorFrame>
}

const StreamResourceEditor = ({
  name,
  onSaveActionChange,
}: {
  name: string
  onSaveActionChange: (action: PipelineNodeResourceSaveAction | null) => void
}) => {
  const { t } = useTranslation()
  const notify = useNotify()

  const [streamConfig, setStreamConfig] = useState<StreamConfigModel | null>(
    null
  )
  const [savedStreamConfig, setSavedStreamConfig] =
    useState<StreamConfigModel | null>(null)
  const [usage, setUsage] = useState(0)
  const dirty = useDirty(savedStreamConfig, streamConfig)

  const [onFetch] = useApiOperation(async (stream: string) => {
    const config = await configApi.getStreamConfig(stream)
    setStreamConfig(config)
    setSavedStreamConfig(config)
    setUsage(await logApi.getStreamUsage(stream))
  }, [])

  useEffect(() => {
    setStreamConfig(null)
    setSavedStreamConfig(null)
    onFetch(name)
  }, [name])

  const [onSave, saveLoading] = useApiOperation(async () => {
    if (streamConfig === null) return
    await configApi.configureStream(name, streamConfig)
    setSavedStreamConfig(streamConfig)
    notify.success(t('pages.storage.notifications.saved'))
  }, [name, streamConfig])

  if (streamConfig === null) {
    return <Loading />
  }

  return (
    <ResourceForm
      saving={saveLoading}
      canSave={dirty}
      onSave={onSave}
      onSaveActionChange={onSaveActionChange}
    >
      <StreamEditor
        stacked
        streamConfig={streamConfig}
        storageUsage={usage}
        onStreamConfigChange={setStreamConfig}
      />
    </ResourceForm>
  )
}

const TransformerResourceEditor = ({
  name,
  onSaveActionChange,
}: {
  name: string
  onSaveActionChange: (action: PipelineNodeResourceSaveAction | null) => void
}) => {
  const { t } = useTranslation()
  const notify = useNotify()

  const [code, setCode] = useState<string | null>(null)
  const [savedCode, setSavedCode] = useState<string | null>(null)
  const dirty = useDirty(savedCode, code)

  const [onFetch] = useApiOperation(async (transformer: string) => {
    const script = await configApi.getTransformer(transformer)
    setCode(script)
    setSavedCode(script)
  }, [])

  useEffect(() => {
    setCode(null)
    setSavedCode(null)
    onFetch(name)
  }, [name])

  const [onSave, saveLoading] = useApiOperation(async () => {
    if (code === null) return
    await configApi.saveTransformer(name, code)
    setSavedCode(code)
    notify.success(t('pages.transformers.notifications.saved'))
  }, [name, code])

  if (code === null) {
    return <Loading />
  }

  return (
    <ResourceForm
      saving={saveLoading}
      canSave={dirty}
      onSave={onSave}
      onSaveActionChange={onSaveActionChange}
    >
      <CodeSurface variant="outlined">
        <VrlCodeEditor
          id="monaco:pipelines.drawer.transformer"
          code={code}
          onCodeChange={setCode}
        />
      </CodeSurface>
    </ResourceForm>
  )
}

const ForwarderResourceEditor = ({
  name,
  onSaveActionChange,
}: {
  name: string
  onSaveActionChange: (action: PipelineNodeResourceSaveAction | null) => void
}) => {
  const { t } = useTranslation()
  const notify = useNotify()

  const [valid, setValid] = useState(false)
  const [forwarder, setForwarder] = useState<ForwarderModel | null>(null)
  const [savedForwarder, setSavedForwarder] = useState<ForwarderModel | null>(
    null
  )
  const dirty = useDirty(savedForwarder, forwarder)

  const initializedRef = useRef(false)
  const handleForwarderChange = (newForwarder: ForwarderModel) => {
    setForwarder(newForwarder)
    if (!initializedRef.current) {
      initializedRef.current = true
      setSavedForwarder(newForwarder)
    }
  }

  const [onFetch] = useApiOperation(async (forwarderName: string) => {
    initializedRef.current = false
    const fetched = await configApi.getForwarder(forwarderName)
    setForwarder(fetched)
    setSavedForwarder(fetched)
  }, [])

  useEffect(() => {
    setForwarder(null)
    setSavedForwarder(null)
    onFetch(name)
  }, [name])

  const [onSave, saveLoading] = useApiOperation(async () => {
    if (forwarder === null) return
    await configApi.saveForwarder(name, forwarder)
    setSavedForwarder(forwarder)
    notify.success(t('pages.forwarders.notifications.saved'))
  }, [name, forwarder])

  if (forwarder === null) {
    return <Loading />
  }

  return (
    <ResourceForm
      saving={saveLoading}
      canSave={dirty && valid}
      onSave={onSave}
      onSaveActionChange={onSaveActionChange}
    >
      <ForwarderEditor
        forwarder={forwarder}
        onForwarderChange={handleForwarderChange}
        onValidationChange={setValid}
      />
    </ResourceForm>
  )
}

const PipelineNodeResourceEditor = ({
  kind,
  name,
  onSaveActionChange,
}: PipelineNodeResourceEditorProps) => {
  switch (kind) {
    case 'stream':
      return (
        <StreamResourceEditor
          key={name}
          name={name}
          onSaveActionChange={onSaveActionChange}
        />
      )
    case 'transformer':
      return (
        <TransformerResourceEditor
          key={name}
          name={name}
          onSaveActionChange={onSaveActionChange}
        />
      )
    case 'forwarder':
      return (
        <ForwarderResourceEditor
          key={name}
          name={name}
          onSaveActionChange={onSaveActionChange}
        />
      )
  }
}

export default PipelineNodeResourceEditor
