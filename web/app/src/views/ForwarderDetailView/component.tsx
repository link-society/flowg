import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { LoaderFunction, useLoaderData, useNavigate } from 'react-router'

import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import Typography from '@mui/material/Typography'

import CancelIcon from '@mui/icons-material/Cancel'
import DeleteIcon from '@mui/icons-material/Delete'
import HelpIcon from '@mui/icons-material/Help'
import PlayArrowIcon from '@mui/icons-material/PlayArrow'
import SaveIcon from '@mui/icons-material/Save'
import ScienceIcon from '@mui/icons-material/Science'

import * as configApi from '@/lib/api/operations/config'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'
import { useProfile } from '@/lib/hooks/profile'

import { loginRequired } from '@/lib/decorators/loaders'

import ButtonNewForwarder from '@/components/ButtonNewForwarder/component'
import ForwarderEditor from '@/components/ForwarderEditor/component'
import InputKeyValue from '@/components/InputKeyValue/component'
import SideNavList from '@/components/SideNavList/component'

import { buildUrl } from '@/router'

import {
  ForwarderDetailViewBody,
  ForwarderDetailViewContent,
  ForwarderDetailViewEditorPaper,
  ForwarderDetailViewHeader,
  ForwarderDetailViewHeaderActions,
  ForwarderDetailViewHeaderLeft,
  ForwarderDetailViewHeaderRight,
  ForwarderDetailViewHeaderTest,
  ForwarderDetailViewRoot,
  ForwarderDetailViewSidebar,
  TestDialogHint,
} from './styles'
import { LoaderData } from './types'

export const loader: LoaderFunction = loginRequired(
  async ({ params }): Promise<LoaderData> => {
    const forwarders = await configApi.listForwarders()

    if (!params.forwarder || !forwarders.includes(params.forwarder)) {
      throw new Response(`Forwarder ${params.forwarder} not found`, {
        status: 404,
      })
    }

    const forwarder = await configApi.getForwarder(params.forwarder)
    return {
      forwarders: forwarders,
      currentForwarder: {
        name: params.forwarder,
        forwarder,
      },
    }
  }
)

const ForwarderDetailView = () => {
  const { t } = useTranslation()
  const notify = useNotify()

  const { permissions } = useProfile()
  const { forwarders, currentForwarder } = useLoaderData() as LoaderData
  const navigate = useNavigate()

  const [forwarder, setForwarder] = useState(currentForwarder.forwarder)
  const [valid, setValid] = useState(false)

  const [testOpen, setTestOpen] = useState(false)
  const [testRecords, setTestRecords] = useState<[string, string][]>([])

  const [onTest, testLoading] = useApiOperation(async () => {
    const input = Object.fromEntries(testRecords)
    await configApi.testForwarder(currentForwarder.name, input)

    notify.success(t('common.notifications.testPassed'))

    setTestOpen(false)
  }, [testRecords])

  const onCreate = (name: string) => {
    queueMicrotask(() => {
      navigate(buildUrl(`/forwarders/${name}`))
    })
  }

  const [onDelete, deleteLoading] = useApiOperation(async () => {
    await configApi.deleteForwarder(currentForwarder.name)
    queueMicrotask(() => {
      navigate(buildUrl('/forwarders'))
    })
  }, [currentForwarder])

  const [onSave, saveLoading] = useApiOperation(async () => {
    await configApi.saveForwarder(currentForwarder.name, forwarder)
    notify.success(t('pages.forwarders.notifications.saved'))
  }, [forwarder, currentForwarder])

  return (
    <>
      <ForwarderDetailViewRoot>
        <ForwarderDetailViewHeader variant="toolbar">
          <ForwarderDetailViewHeaderLeft>
            <Button
              variant="contained"
              color="primary"
              size="small"
              href="https://flowg.cloud/docs/"
              target="_blank"
              startIcon={<HelpIcon />}
            >
              {t('common.actions.documentation')}
            </Button>
          </ForwarderDetailViewHeaderLeft>

          <ForwarderDetailViewHeaderRight>
            <ForwarderDetailViewHeaderTest>
              <Button
                variant="contained"
                color="primary"
                size="small"
                onClick={() => setTestOpen(true)}
                startIcon={<ScienceIcon />}
              >
                {t('common.actions.test')}
              </Button>
            </ForwarderDetailViewHeaderTest>

            {permissions.can_edit_forwarders && (
              <ForwarderDetailViewHeaderActions>
                <ButtonNewForwarder onForwarderCreated={onCreate} />

                <Button
                  id="btn:forwarders.delete"
                  variant="contained"
                  color="error"
                  size="small"
                  onClick={onDelete}
                  disabled={deleteLoading}
                  startIcon={!deleteLoading && <DeleteIcon />}
                >
                  {deleteLoading ? (
                    <CircularProgress size={24} />
                  ) : (
                    <>{t('common.actions.delete')}</>
                  )}
                </Button>

                <Button
                  id="btn:forwarders.save"
                  variant="contained"
                  color="secondary"
                  size="small"
                  onClick={onSave}
                  disabled={saveLoading || !valid}
                  startIcon={!saveLoading && <SaveIcon />}
                >
                  {saveLoading ? (
                    <CircularProgress size={24} />
                  ) : (
                    <>{t('common.actions.save')}</>
                  )}
                </Button>
              </ForwarderDetailViewHeaderActions>
            )}
          </ForwarderDetailViewHeaderRight>
        </ForwarderDetailViewHeader>

        <ForwarderDetailViewBody variant="page">
          <ForwarderDetailViewSidebar>
            <SideNavList
              namespace="forwarders"
              urlPrefix={buildUrl('/forwarders')}
              items={forwarders}
              currentItem={currentForwarder.name}
            />
          </ForwarderDetailViewSidebar>
          <ForwarderDetailViewContent>
            <ForwarderDetailViewEditorPaper>
              <ForwarderEditor
                forwarder={forwarder}
                onForwarderChange={setForwarder}
                onValidationChange={setValid}
              />
            </ForwarderDetailViewEditorPaper>
          </ForwarderDetailViewContent>
        </ForwarderDetailViewBody>
      </ForwarderDetailViewRoot>

      <Dialog open={testOpen} scroll="paper" onClose={() => setTestOpen(false)}>
        <DialogTitle>
          <Typography variant="titleMd">
            {t('common.testDialog.title')}
          </Typography>
        </DialogTitle>
        <DialogContent>
          <TestDialogHint>
            <Typography variant="text">
              {t('common.testDialog.inputRecord')}
            </Typography>
          </TestDialogHint>
          <InputKeyValue
            id="kv:transformers.test.record"
            keyLabel={t('common.testDialog.fieldLabel')}
            keyValues={testRecords}
            onChange={setTestRecords}
          />
        </DialogContent>
        <DialogActions>
          <Button
            variant="contained"
            startIcon={<CancelIcon />}
            onClick={() => setTestOpen(false)}
            disabled={testLoading}
          >
            {t('common.actions.cancel')}
          </Button>
          <Button
            id="btn:transformers.test.run"
            variant="contained"
            color="secondary"
            endIcon={<PlayArrowIcon />}
            disabled={testLoading}
            onClick={() => onTest()}
          >
            {testLoading ? (
              <CircularProgress size={24} />
            ) : (
              <>{t('common.actions.run')}</>
            )}
          </Button>
        </DialogActions>
      </Dialog>
    </>
  )
}

export default ForwarderDetailView
