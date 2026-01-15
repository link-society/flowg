import { useEffect } from 'react'
import { LoaderFunction, useLoaderData, useNavigate } from 'react-router'

import * as configApi from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

import ButtonNewForwarder from '@/components/ButtonNewForwarder'

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
    <div className="w-full h-full flex flex-col items-center justify-center gap-5">
      <h1 className="text-3xl font-semibold">No forwarder found, create one</h1>

      <ButtonNewForwarder
        onForwarderCreated={(name) => {
          navigate(`/web/forwarders/${name}`)
        }}
      />
    </div>
  )
}

export default ForwarderSectionView
