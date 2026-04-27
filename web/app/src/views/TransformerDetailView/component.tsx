import { useState } from 'react'
import { LoaderFunction, useLoaderData, useNavigate } from 'react-router'

import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'

import DeleteIcon from '@mui/icons-material/Delete'
import HelpIcon from '@mui/icons-material/Help'
import SaveIcon from '@mui/icons-material/Save'

import * as configApi from '@/lib/api/operations/config'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'
import { useProfile } from '@/lib/hooks/profile'

import { loginRequired } from '@/lib/decorators/loaders'

import ButtonNewTransformer from '@/components/ButtonNewTransformer/component'
import SideNavList from '@/components/SideNavList/component'
import TransformerEditor from '@/components/TransformerEditor/component'

import {
  TransformerDetailViewBody,
  TransformerDetailViewContainer,
  TransformerDetailViewEditor,
  TransformerDetailViewSidebar,
  TransformerDetailViewToolbar,
  TransformerDetailViewToolbarLeft,
  TransformerDetailViewToolbarRight,
} from './styles'
import { LoaderData } from './types'

export const loader: LoaderFunction = loginRequired(
  async ({ params }): Promise<LoaderData> => {
    const transformers = await configApi.listTransformers()

    if (!transformers.includes(params.transformer!)) {
      throw new Response(`Transformer ${params.transformer} not found`, {
        status: 404,
      })
    }

    const script = await configApi.getTransformer(params.transformer!)
    return {
      transformers,
      currentTransformer: {
        name: params.transformer!,
        script,
      },
    }
  }
)

const TransformerDetailView = () => {
  const notify = useNotify()

  const { permissions } = useProfile()
  const { transformers, currentTransformer } = useLoaderData() as LoaderData
  const navigate = useNavigate()

  const [code, setCode] = useState(currentTransformer.script)

  const onCreate = (name: string) => {
    queueMicrotask(() => {
      navigate(`/web/transformers/${name}`)
    })
  }

  const [onDelete, deleteLoading] = useApiOperation(async () => {
    await configApi.deleteTransformer(currentTransformer.name)
    queueMicrotask(() => {
      navigate('/web/transformers')
    })
  }, [currentTransformer])

  const [onSave, saveLoading] = useApiOperation(async () => {
    await configApi.saveTransformer(currentTransformer.name, code)
    notify.success('Transformer saved')
  }, [code, currentTransformer])

  return (
    <TransformerDetailViewContainer>
      <TransformerDetailViewToolbar variant="toolbar">
        <TransformerDetailViewToolbarLeft>
          <Button
            variant="contained"
            color="primary"
            size="small"
            href="https://vector.dev/docs/reference/vrl/"
            target="_blank"
            startIcon={<HelpIcon />}
          >
            VRL Documentation
          </Button>
        </TransformerDetailViewToolbarLeft>

        {permissions.can_edit_transformers && (
          <TransformerDetailViewToolbarRight>
            <ButtonNewTransformer onTransformerCreated={onCreate} />

            <Button
              id="btn:transformers.delete"
              variant="contained"
              color="error"
              size="small"
              onClick={onDelete}
              disabled={deleteLoading}
              startIcon={!deleteLoading && <DeleteIcon />}
            >
              {deleteLoading ? <CircularProgress size={24} /> : <>Delete</>}
            </Button>

            <Button
              id="btn:transformers.save"
              variant="contained"
              color="secondary"
              size="small"
              onClick={onSave}
              disabled={saveLoading}
              startIcon={!saveLoading && <SaveIcon />}
            >
              {saveLoading ? <CircularProgress size={24} /> : <>Save</>}
            </Button>
          </TransformerDetailViewToolbarRight>
        )}
      </TransformerDetailViewToolbar>

      <TransformerDetailViewBody variant="page">
        <TransformerDetailViewSidebar>
          <SideNavList
            namespace="transformers"
            urlPrefix="/web/transformers"
            items={transformers}
            currentItem={currentTransformer.name}
          />
        </TransformerDetailViewSidebar>

        <TransformerDetailViewEditor>
          <TransformerEditor code={code} onCodeChange={setCode} />
        </TransformerDetailViewEditor>
      </TransformerDetailViewBody>
    </TransformerDetailViewContainer>
  )
}

export default TransformerDetailView
