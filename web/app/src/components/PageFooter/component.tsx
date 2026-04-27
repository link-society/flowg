import { Typography } from '@mui/material'

import { useFeatureFlags } from '@/lib/hooks/featureflags'

import { PageFooterContainer } from './styles'

const PageFooter = () => {
  const { demoMode } = useFeatureFlags()

  return (
    <PageFooterContainer>
      {demoMode && (
        <Typography variant="text" fontWeight={600}>
          Demo Mode Enabled
        </Typography>
      )}

      <div>
        <Typography variant="text" fontWeight={600}>
          {import.meta.env.FLOWG_VERSION}
        </Typography>
      </div>
    </PageFooterContainer>
  )
}

export default PageFooter
