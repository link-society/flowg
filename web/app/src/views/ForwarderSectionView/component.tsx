import { LoaderFunction, redirect, useNavigate } from 'react-router'

import Typography from '@mui/material/Typography'

import * as configApi from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

import ButtonNewForwarder from '@/components/ButtonNewForwarder/component'

import { ForwarderSectionViewRoot } from './styles'

export const loader: LoaderFunction = loginRequired(async () => {
  const forwarders = await configApi.listForwarders()
  if (forwarders.length > 0) {
    return redirect(`/web/forwarders/${forwarders[0]}`)
  }

  return null
})

const ForwarderSectionView = () => {
  const navigate = useNavigate()

  return (
    <ForwarderSectionViewRoot>
      <Typography variant="titleLg" component="h1">
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
