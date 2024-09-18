import { useEffect } from 'react'
import { useLoaderData, useNavigate } from 'react-router-dom'

import Backdrop from '@mui/material/Backdrop'
import CircularProgress from '@mui/material/CircularProgress'

export const TransformerView = () => {
  const { transformers } = useLoaderData() as { transformers: string[] }
  const navigate = useNavigate()

  useEffect(
    () => {
      if (transformers.length > 0) {
        navigate(`/web/transformers/${transformers[0]}`)
      }
    },
    [],
  )

  return (
    <>
      {transformers.length > 0
        ? (
          <Backdrop open={true}>
            <CircularProgress color="inherit" />
          </Backdrop>
        )
        : (
          <div>No data</div>
        )
      }
    </>
  )
}
