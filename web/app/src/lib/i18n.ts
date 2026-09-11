import dummyTranslationsUrl from '@/locales/generated/dummy/translation.json?url&no-inline'
import enTranslationsUrl from '@/locales/generated/en/translation.json?url&no-inline'
import i18n from 'i18next'
import HttpBackend from 'i18next-http-backend'

import { initReactI18next } from 'react-i18next'

const translations: Record<string, string> = {
  dummy: dummyTranslationsUrl,
  en: enTranslationsUrl,
}

export type Language = {
  code: string
  label: string
}

export const AVAILABLE_LANGUAGES: Language[] = [
  { code: 'en', label: 'English' },
]

export const LANGUAGE_STORAGE_KEY = 'lng'

const availableCodes = AVAILABLE_LANGUAGES.map((language) => language.code)
const storedLanguage = localStorage.getItem(LANGUAGE_STORAGE_KEY)
const initialLanguage =
  storedLanguage && availableCodes.includes(storedLanguage)
    ? storedLanguage
    : 'en'

i18n
  .use(HttpBackend)
  .use(initReactI18next)
  .init({
    lng: initialLanguage,
    fallbackLng: ['dummy', 'en'],
    supportedLngs: ['en', 'dummy'],
    load: 'languageOnly',
    returnEmptyString: false,
    backend: {
      loadPath: (languages: string[], namespaces: string[]) => {
        const language = languages[0]
        const namespace = namespaces[0]

        if (namespace !== 'translation') {
          return false
        }

        return translations[language]
      },
    },
    interpolation: { escapeValue: false },
    react: { useSuspense: true },
  })

export default i18n
