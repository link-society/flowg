import { useEffect } from 'react'
import { LoaderFunction, useLoaderData, useNavigate } from 'react-router'

import Typography from '@mui/material/Typography'

import * as configApi from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

import ButtonNewForwarder from '@/components/ButtonNewForwarder/component'

import { ForwarderSectionViewRoot } from './styles'

export const loader: LoaderFunction = loginRequired(async () => {
  const forwarders = await configApi.listForwarders()
  if (forwarders.length > 0) {
    return { redirectTo: `/web/forwarders/${forwarders[0]}` }
  }

  return { redirectTo: null }
})

const ForwarderSectionView = () => {
  const navigate = useNavigate()
  const { redirectTo } = useLoaderData<{ redirectTo: string | null }>()

  useEffect(() => {
    if (redirectTo !== null) {
      navigate(redirectTo, { replace: true })
    }
  }, [redirectTo])

  return (
    <ForwarderSectionViewRoot>
      <Typography variant="titleLg" fontWeight={700} component="h1">
        No forwarder found, create one
      </Typography>

      <ButtonNewForwarder
        onForwarderCreated={(name) => {
          navigate(`/web/forwarders/${name}`)
        }}
      />
    </ForwarderSectionViewRoot>
  )
}

export default ForwarderSectionView
