import { useFeatureFlags } from '@/lib/hooks/featureflags'

import { PageFooterContainer, PageFooterText } from './styles'

const PageFooter = () => {
  const { demoMode } = useFeatureFlags()

  return (
    <PageFooterContainer>
      {demoMode && (
        <PageFooterText variant="text">Demo Mode Enabled</PageFooterText>
      )}

      <div>
        <PageFooterText variant="text">
          {import.meta.env.FLOWG_VERSION}
        </PageFooterText>
      </div>
    </PageFooterContainer>
  )
}

export default PageFooter
