import i18n from 'i18next'
import HttpBackend from 'i18next-http-backend'

import { initReactI18next } from 'react-i18next'

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
    backend: { loadPath: './assets/locales/{{lng}}/{{ns}}.json' },
    interpolation: { escapeValue: false },
    react: { useSuspense: true },
  })

export default i18n
