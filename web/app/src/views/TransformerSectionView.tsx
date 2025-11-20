import { LoaderFunction, redirect, useNavigate } from 'react-router'

import * as configApi from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

import ButtonNewTransformer from '@/components/ButtonNewTransformer'

export const loader: LoaderFunction = loginRequired(async () => {
  const transformers = await configApi.listTransformers()
  if (transformers.length > 0) {
    throw redirect(`/web/transformers/${transformers[0]}`)
  }
})

const TransformerSectionView = () => {
  const navigate = useNavigate()

  return (
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
  )
}

export default TransformerSectionView
