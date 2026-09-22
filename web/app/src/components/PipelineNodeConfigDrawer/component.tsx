import { ReactNode, useCallback, useEffect, useState } from 'react'
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

type InlineFieldEditorProps = Readonly<{
  nodeId: string
  dataKey: string
  label: string
  value: string
}>

/** Editable single-field body for `switch` / `metric` nodes. */
const InlineFieldEditor = ({
  nodeId,
  dataKey,
  label,
  value,
}: InlineFieldEditorProps) => {
  const { setNodes } = usePipelineEditorHooks()
  const [current, setCurrent] = useState(value)

  useEffect(() => {
    setNodes((prevNodes) =>
      prevNodes.map((node) =>
        node.id === nodeId
          ? { ...node, data: { ...node.data, [dataKey]: current } }
          : node
      )
    )
  }, [nodeId, dataKey, current, setNodes])

  return (
    <TextField
      label={label}
      type="text"
      value={current}
      onChange={(evt) => setCurrent(evt.target.value)}
      fullWidth
      variant="outlined"
      slotProps={{ input: { sx: { fontFamily: 'monospace' } } }}
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
              <InlineFieldEditor
                nodeId={node.id}
                dataKey="condition"
                label={t('components.pipelineNodeSwitch.label')}
                value={str(data.condition)}
              />
            )}

            {nodeType === 'metric' && (
              <InlineFieldEditor
                nodeId={node.id}
                dataKey="name"
                label={t('components.pipelineNodeMetric.label')}
                value={str(data.name)}
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
