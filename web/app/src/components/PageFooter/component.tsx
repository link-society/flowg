import { useTranslation } from 'react-i18next'

import Divider from '@mui/material/Divider'

import { useFeatureFlags } from '@/lib/hooks/featureflags'

import ColorSelect from '@/components/ColorSelect/component'
import LanguageSelect from '@/components/LanguageSelect/component'

import {
  PageFooterActions,
  PageFooterContainer,
  PageFooterText,
} from './styles'

const PageFooter = () => {
  const { t } = useTranslation()
  const { demoMode } = useFeatureFlags()

  return (
    <PageFooterContainer>
      <PageFooterActions>
        <ColorSelect />

        <Divider orientation="vertical" flexItem />

        <LanguageSelect />
      </PageFooterActions>

      {demoMode && (
        <PageFooterText variant="text">
          {t('components.pageFooter.demoModeEnabled')}
        </PageFooterText>
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
