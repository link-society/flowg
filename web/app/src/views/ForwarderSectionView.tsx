import { useEffect } from 'react'
import { LoaderFunction, useLoaderData, useNavigate } from 'react-router'

import Backdrop from '@mui/material/Backdrop'
import CircularProgress from '@mui/material/CircularProgress'

import * as configApi from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

import ButtonNewForwarder from '@/components/ButtonNewForwarder'

export type LoaderData = {
  forwarders: string[]
}

export const loader: LoaderFunction = loginRequired(
  async (): Promise<LoaderData> => {
    const forwarders = await configApi.listForwarders()
    return { forwarders }
  }
)

const ForwarderSectionView = () => {
  const { forwarders } = useLoaderData() as LoaderData
  const navigate = useNavigate()

  useEffect(() => {
    if (forwarders.length > 0) {
      navigate(`/web/forwarders/${forwarders[0]}`, { replace: true })
    }
  }, [])

  return (
    <>
      {forwarders.length > 0 ? (
        <Backdrop open={true}>
          <CircularProgress color="inherit" />
        </Backdrop>
      ) : (
        <div className="w-full h-full flex flex-col items-center justify-center gap-5">
          <h1 className="text-3xl font-semibold">
            No forwarder found, create one
          </h1>

          <ButtonNewForwarder
            onForwarderCreated={(name) => {
              navigate(`/web/forwarders/${name}`)
            }}
          />
        </div>
      )}
    </>
  )
}

export default ForwarderSectionView
