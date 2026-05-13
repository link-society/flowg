import { Typography } from '@mui/material'

import { LoaderFunction, redirect, useNavigate } from 'react-router'

import * as configApi from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

import ButtonNewTransformer from '@/components/ButtonNewTransformer/component'

import { buildUrl } from '@/router'

import { TransformerSectionViewContainer } from './styles'

export const loader: LoaderFunction = loginRequired(async () => {
  const transformers = await configApi.listTransformers()
  if (transformers.length > 0) {
    return redirect(buildUrl(`/transformers/${transformers[0]}`))
  }

  return null
})

const TransformerSectionView = () => {
  const navigate = useNavigate()

  return (
    <TransformerSectionViewContainer>
      <Typography variant="titleLg" component="h1">
        No transformer found, create one
      </Typography>

      <ButtonNewTransformer
        onTransformerCreated={(name) => {
          navigate(buildUrl(`/transformers/${name}`))
        }}
      />
    </TransformerSectionViewContainer>
  )
}

export default TransformerSectionView
