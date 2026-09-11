import { ReactNode, useCallback, useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import IconButton from '@mui/material/IconButton'
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
import StorageIcon from '@mui/icons-material/Storage'

import { Node, useOnSelectionChange } from '@xyflow/react'

import { usePipelineEditorHooks } from '@/lib/hooks/pipeline-editor'
import { useProfile } from '@/lib/hooks/profile'

import DialogForwarderEditor from '@/components/DialogForwarderEditor/component'
import DialogStreamEditor from '@/components/DialogStreamEditor/component'
import DialogTransformerEditor from '@/components/DialogTransformerEditor/component'
import PipelineDeleteNodeButton from '@/components/PipelineDeleteNodeButton/component'

import {
  DrawerAccent,
  DrawerActions,
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

  const [selectedNode, setSelectedNode] = useState<Node | null>(null)
  const [closedNodeId, setClosedNodeId] = useState<string | null>(null)

  const onSelectionChange = useCallback(({ nodes }: { nodes: Node[] }) => {
    setSelectedNode(nodes.length === 1 ? nodes[0] : null)
  }, [])
  useOnSelectionChange({ onChange: onSelectionChange })

  if (selectedNode === null || selectedNode.id === closedNodeId) {
    return null
  }

  const node = selectedNode
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
    <DrawerRoot>
      <DrawerAccent style={{ backgroundColor: accent }} />

      <DrawerHeader>
        <DrawerHeaderIcon style={{ backgroundColor: accent }}>
          {meta.icon}
        </DrawerHeaderIcon>
        <DrawerHeaderText>
          <DrawerHeaderKicker>{t(meta.labelKey)}</DrawerHeaderKicker>
          <DrawerHeaderTitle title={title}>{title}</DrawerHeaderTitle>
        </DrawerHeaderText>
        <IconButton
          size="small"
          edge="end"
          sx={{ marginLeft: 'auto' }}
          onClick={() => setClosedNodeId(node.id)}
          aria-label={t('common.actions.close')}
        >
          <CloseIcon fontSize="small" />
        </IconButton>
      </DrawerHeader>

      <DrawerBody key={node.id}>
        {nodeType === 'source' && (
          <>
            <TextField
              label={t('components.pipelineNodeConfigDrawer.sourceTypeLabel')}
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
              {t('components.pipelineNodeConfigDrawer.sharedResourceHint')}
            </DrawerHint>
            {permissions.can_edit_transformers && (
              <DrawerActions>
                <DialogTransformerEditor transformer={str(data.transformer)} />
              </DrawerActions>
            )}
          </>
        )}

        {nodeType === 'forwarder' && (
          <>
            <DrawerHint>
              {t('components.pipelineNodeConfigDrawer.sharedResourceHint')}
            </DrawerHint>
            {permissions.can_edit_forwarders && (
              <DrawerActions>
                <DialogForwarderEditor forwarderName={str(data.forwarder)} />
              </DrawerActions>
            )}
          </>
        )}

        {nodeType === 'router' && (
          <>
            <DrawerHint>
              {t('components.pipelineNodeConfigDrawer.sharedResourceHint')}
            </DrawerHint>
            {permissions.can_edit_streams && (
              <DrawerActions>
                <DialogStreamEditor stream={str(data.stream)} />
              </DrawerActions>
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

        {canDelete && (
          <DrawerActions>
            <PipelineDeleteNodeButton nodeId={node.id} />
          </DrawerActions>
        )}
      </DrawerBody>

      <DrawerFooter>
        <InfoOutlinedIcon
          fontSize="small"
          sx={{ color: theme.tokens.colors.mutedText, flex: 'none' }}
        />
        <span>{footerNote}</span>
      </DrawerFooter>
    </DrawerRoot>
  )
}

export default PipelineNodeConfigDrawer
