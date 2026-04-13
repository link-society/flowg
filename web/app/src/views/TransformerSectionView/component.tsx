import { Typography } from '@mui/material'

import { useEffect } from 'react'
import { LoaderFunction, useLoaderData, useNavigate } from 'react-router'

import * as configApi from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

import ButtonNewTransformer from '@/components/ButtonNewTransformer/component'

import { TransformerSectionViewContainer } from './styles'

export const loader: LoaderFunction = loginRequired(async () => {
  const transformers = await configApi.listTransformers()
  if (transformers.length > 0) {
    return { redirectTo: `/web/transformers/${transformers[0]}` }
  }

  return { redirectTo: null }
})

const TransformerSectionView = () => {
  const navigate = useNavigate()
  const { redirectTo } = useLoaderData<{ redirectTo: string | null }>()

  useEffect(() => {
    if (redirectTo !== null) {
      navigate(redirectTo, { replace: true })
    }
  }, [redirectTo])

  return (
    <TransformerSectionViewContainer>
      <Typography variant="titleLg" component="h1" fontWeight={600}>
        No transformer found, create one
      </Typography>

      <ButtonNewTransformer
        onTransformerCreated={(name) => {
          navigate(`/web/transformers/${name}`)
        }}
      />
    </TransformerSectionViewContainer>
  )
}

export default TransformerSectionView
