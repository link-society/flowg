import { ChangeEventHandler, ReactNode, useCallback, useEffect, useState } from 'react'
import { createPortal } from 'react-dom'
import { useTranslation } from 'react-i18next'

import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import Fade from '@mui/material/Fade'
import IconButton from '@mui/material/IconButton'
import Slide from '@mui/material/Slide'
import TextField from '@mui/material/TextField'
import { useTheme } from '@mui/material/styles'

import AccountTreeIcon from '@mui/icons-material/AccountTree'
import BarChartIcon from '@mui/icons-material/BarChart'
import CloseIcon from '@mui/icons-material/Close'
import DeviceHubIcon from '@mui/icons-material/DeviceHub'
import FilterAltIcon from '@mui/icons-material/FilterAlt'
import ForwardToInboxIcon from '@mui/icons-material/ForwardToInbox'
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined'
import InputIcon from '@mui/icons-material/Input'
import SaveIcon from '@mui/icons-material/Save'
import StorageIcon from '@mui/icons-material/Storage'

import { Node, useOnSelectionChange } from '@xyflow/react'

import { usePipelineEditorHooks } from '@/lib/hooks/pipeline-editor'
import { usePipelineOverlayHost } from '@/lib/hooks/pipeline-overlay-host'
import { useProfile } from '@/lib/hooks/profile'

import PipelineDeleteNodeButton from '@/components/PipelineDeleteNodeButton/component'
import PipelineNodeResourceEditor from '@/components/PipelineNodeResourceEditor/component'
import { PipelineNodeResourceSaveAction } from '@/components/PipelineNodeResourceEditor/types'

import {
  DrawerAccent,
  DrawerActions,
  DrawerBackdrop,
  DrawerBody,
  DrawerFooter,
  DrawerHeader,
  DrawerHeaderIcon,
  DrawerHeaderKicker,
  DrawerHeaderText,
  DrawerHeaderTitle,
  DrawerHint,
  DrawerRoot,
} from './styles'

const str = (value: unknown): string => (typeof value === 'string' ? value : '')

type NodeTypeMeta = {
  icon: ReactNode
  labelKey: string
  accentToken:
    | 'nodeSourceBg'
    | 'nodeTransformerBg'
    | 'nodeForwarderBg'
    | 'nodeRouterBg'
    | 'nodeSwitchBg'
    | 'nodeMetricBg'
    | 'nodePipelineBg'
}

const NODE_META: Record<string, NodeTypeMeta> = {
  source: {
    icon: <InputIcon fontSize="small" />,
    labelKey: 'components.pipelineNodeSource.label',
    accentToken: 'nodeSourceBg',
  },
  transform: {
    icon: <FilterAltIcon fontSize="small" />,
    labelKey: 'components.pipelineNodeTransformer.label',
    accentToken: 'nodeTransformerBg',
  },
  forwarder: {
    icon: <ForwardToInboxIcon fontSize="small" />,
    labelKey: 'components.pipelineNodeForwarder.label',
    accentToken: 'nodeForwarderBg',
  },
  router: {
    icon: <StorageIcon fontSize="small" />,
    labelKey: 'components.pipelineNodeRouter.label',
    accentToken: 'nodeRouterBg',
  },
  switch: {
    icon: <DeviceHubIcon fontSize="small" />,
    labelKey: 'components.pipelineNodeSwitch.label',
    accentToken: 'nodeSwitchBg',
  },
  metric: {
    icon: <BarChartIcon fontSize="small" />,
    labelKey: 'components.pipelineNodeMetric.label',
    accentToken: 'nodeMetricBg',
  },
  pipeline: {
    icon: <AccountTreeIcon fontSize="small" />,
    labelKey: 'components.pipelineNodePipeline.label',
    accentToken: 'nodePipelineBg',
  },
}

const SHARED_RESOURCE_TYPES = new Set(['transform', 'forwarder', 'router'])

const NodeDataField = ({
  nodeId,
  field,
  label,
  initialValue,
  monospace,
}: {
  nodeId: string
  field: string
  label: string
  initialValue: string
  monospace?: boolean
}) => {
  const { setNodes } = usePipelineEditorHooks()
  const [value, setValue] = useState(initialValue)

  const onChange: ChangeEventHandler<HTMLInputElement> = (evt) => {
    const newValue = evt.target.value
    setValue(newValue)
    setNodes((nds) =>
      nds.map((n) =>
        n.id === nodeId ? { ...n, data: { ...n.data, [field]: newValue } } : n
      )
    )
  }

  return (
    <TextField
      label={label}
      type="text"
      value={value}
      onChange={onChange}
      fullWidth
      variant="outlined"
      slotProps={{
        input: { sx: monospace ? { fontFamily: 'monospace' } : undefined },
      }}
    />
  )
}

const PipelineNodeConfigDrawer = () => {
  const { t } = useTranslation()
  const theme = useTheme()
  const { permissions } = useProfile()
  const overlayHost = usePipelineOverlayHost()
  const { setNodes } = usePipelineEditorHooks()

  const [selectedNode, setSelectedNode] = useState<Node | null>(null)
  const [displayedNode, setDisplayedNode] = useState<Node | null>(null)
  const [saveAction, setSaveAction] =
    useState<PipelineNodeResourceSaveAction | null>(null)

  const onSelectionChange = useCallback(({ nodes }: { nodes: Node[] }) => {
    setSelectedNode(nodes.length === 1 ? nodes[0] : null)
  }, [])
  useOnSelectionChange({ onChange: onSelectionChange })

  const handleClose = useCallback(
    (nodeId: string) => {
      setNodes((nds) =>
        nds.map((n) => (n.id === nodeId ? { ...n, selected: false } : n))
      )
    },
    [setNodes]
  )

  const open = selectedNode !== null

  useEffect(() => {
    if (open) {
      setDisplayedNode(selectedNode)
    }
  }, [open, selectedNode])

  if (displayedNode === null) {
    return null
  }

  const node = displayedNode
  const nodeType = node.type ?? ''
  const meta = NODE_META[nodeType]
  if (meta === undefined) {
    return null
  }

  const data = node.data as Record<string, unknown>
  const accent = theme.tokens.colors[meta.accentToken]

  let title: string
  switch (nodeType) {
    case 'source':
      title = str(data.type).toUpperCase()
      break
    case 'transform':
      title = str(data.transformer)
      break
    case 'forwarder':
      title = str(data.forwarder)
      break
    case 'router':
      title = str(data.stream)
      break
    case 'pipeline':
      title = str(data.pipeline)
      break
    case 'switch':
      title = t('components.pipelineEditorFlow.switchNodeLabel')
      break
    case 'metric':
      title = t('components.pipelineEditorFlow.metricNodeLabel')
      break
    default:
      title = nodeType
  }

  const canDelete = node.deletable !== false && nodeType !== 'source'
  const footerNote = SHARED_RESOURCE_TYPES.has(nodeType)
    ? t('components.pipelineNodeConfigDrawer.sharedResourceNote')
    : t('components.pipelineNodeConfigDrawer.pipelineStoredNote')

  return (
    <>
      {overlayHost !== null &&
        createPortal(
          <Fade in={open} unmountOnExit>
            <DrawerBackdrop onClick={() => handleClose(node.id)} />
          </Fade>,
          overlayHost
        )}
      <Slide direction="left" in={open} onExited={() => setDisplayedNode(null)}>
        <DrawerRoot>
          <DrawerAccent style={{ backgroundColor: accent }} />

          <DrawerHeader>
            <DrawerHeaderIcon style={{ backgroundColor: accent }}>
              {meta.icon}
            </DrawerHeaderIcon>
            <DrawerHeaderText>
              <DrawerHeaderKicker variant="text">
                {t(meta.labelKey)}
              </DrawerHeaderKicker>
              <DrawerHeaderTitle variant="titleSm" title={title}>
                {title}
              </DrawerHeaderTitle>
            </DrawerHeaderText>
            <IconButton
              size="small"
              edge="end"
              sx={{ marginLeft: 'auto' }}
              onClick={() => handleClose(node.id)}
              aria-label={t('common.actions.close')}
            >
              <CloseIcon fontSize="small" />
            </IconButton>
          </DrawerHeader>

          <DrawerBody key={node.id}>
            {nodeType === 'source' && (
              <>
                <TextField
                  label={t(
                    'components.pipelineNodeConfigDrawer.sourceTypeLabel'
                  )}
                  type="text"
                  value={str(data.type).toUpperCase()}
                  fullWidth
                  variant="outlined"
                  slotProps={{ input: { readOnly: true } }}
                />
                <DrawerHint>
                  {t('components.pipelineNodeConfigDrawer.sourcePlaceholder')}
                </DrawerHint>
              </>
            )}

            {nodeType === 'transform' && (
              <>
                <DrawerHint>
                  {t(
                    permissions.can_edit_transformers
                      ? 'components.pipelineNodeConfigDrawer.sharedResourceHint'
                      : 'components.pipelineNodeConfigDrawer.readOnlyHint'
                  )}
                </DrawerHint>
                {permissions.can_edit_transformers && (
                  <PipelineNodeResourceEditor
                    kind="transformer"
                    name={str(data.transformer)}
                    onSaveActionChange={setSaveAction}
                  />
                )}
              </>
            )}

            {nodeType === 'forwarder' && (
              <>
                <DrawerHint>
                  {t(
                    permissions.can_edit_forwarders
                      ? 'components.pipelineNodeConfigDrawer.sharedResourceHint'
                      : 'components.pipelineNodeConfigDrawer.readOnlyHint'
                  )}
                </DrawerHint>
                {permissions.can_edit_forwarders && (
                  <PipelineNodeResourceEditor
                    kind="forwarder"
                    name={str(data.forwarder)}
                    onSaveActionChange={setSaveAction}
                  />
                )}
              </>
            )}

            {nodeType === 'router' && (
              <>
                <DrawerHint>
                  {t(
                    permissions.can_edit_streams
                      ? 'components.pipelineNodeConfigDrawer.sharedResourceHint'
                      : 'components.pipelineNodeConfigDrawer.readOnlyHint'
                  )}
                </DrawerHint>
                {permissions.can_edit_streams && (
                  <PipelineNodeResourceEditor
                    kind="stream"
                    name={str(data.stream)}
                    onSaveActionChange={setSaveAction}
                  />
                )}
              </>
            )}

            {nodeType === 'switch' && (
              <NodeDataField
                nodeId={node.id}
                field="condition"
                label={t('components.pipelineNodeSwitch.label')}
                initialValue={str(data.condition)}
                monospace
              />
            )}

            {nodeType === 'metric' && (
              <NodeDataField
                nodeId={node.id}
                field="name"
                label={t('components.pipelineNodeMetric.label')}
                initialValue={str(data.name)}
                monospace
              />
            )}

            {nodeType === 'pipeline' && (
              <TextField
                label={t('components.pipelineNodePipeline.label')}
                type="text"
                value={str(data.pipeline)}
                fullWidth
                variant="outlined"
                slotProps={{ input: { readOnly: true } }}
              />
            )}
          </DrawerBody>

          {(saveAction !== null || canDelete) && (
            <DrawerActions>
              {canDelete && <PipelineDeleteNodeButton nodeId={node.id} />}
              {saveAction !== null && (
                <Button
                  variant="contained"
                  color="secondary"
                  size="small"
                  onClick={saveAction.onSave}
                  disabled={saveAction.saving || !saveAction.canSave}
                  startIcon={!saveAction.saving && <SaveIcon />}
                >
                  {saveAction.saving ? (
                    <CircularProgress size={20} />
                  ) : (
                    t('common.actions.save')
                  )}
                </Button>
              )}
            </DrawerActions>
          )}

          <DrawerFooter>
            <InfoOutlinedIcon
              fontSize="small"
              sx={{ color: theme.tokens.colors.mutedText, flex: 'none' }}
            />
            <span>{footerNote}</span>
          </DrawerFooter>
        </DrawerRoot>
      </Slide>
    </>
  )
}

export default PipelineNodeConfigDrawer
