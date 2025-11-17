import { useEffect } from 'react'
import { LoaderFunction, useLoaderData, useNavigate } from 'react-router'

import Backdrop from '@mui/material/Backdrop'
import CircularProgress from '@mui/material/CircularProgress'

import * as configApi from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

import ButtonNewTransformer from '@/components/ButtonNewTransformer'

type LoaderData = {
  transformers: string[]
}

export const loader: LoaderFunction = loginRequired(
  async (): Promise<LoaderData> => {
    const transformers = await configApi.listTransformers()
    return { transformers }
  }
)

const TransformerSectionView = () => {
  const { transformers } = useLoaderData() as LoaderData
  const navigate = useNavigate()

  useEffect(() => {
    if (transformers.length > 0) {
      navigate(`/web/transformers/${transformers[0]}`, { replace: true })
    }
  }, [])

  return (
    <>
      {transformers.length > 0 ? (
        <Backdrop open={true}>
          <CircularProgress color="inherit" />
        </Backdrop>
      ) : (
        <div className="w-full h-full flex flex-col items-center justify-center gap-5">
          <h1 className="text-3xl font-semibold">
            No transformer found, create one
          </h1>

          <ButtonNewTransformer
            onTransformerCreated={(name) => {
              navigate(`/web/transformers/${name}`)
            }}
          />
        </div>
      )}
    </>
  )
}

export default TransformerSectionView
