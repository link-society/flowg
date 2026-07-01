import siteConfig from '@generated/docusaurus.config'
// eslint-disable-next-line import/no-extraneous-dependencies
import Prism from 'prismjs'

/**
 * Swizzled override of Docusaurus' `prism-include-languages`.
 *
 * `redocusaurus`/`redoc` loads `prismjs/components/*` modules first, attaching
 * their grammars to the standalone `prismjs` instance and populating the module
 * cache. Docusaurus' default implementation then re-`require`s those same
 * component modules to attach them to prism-react-renderer's *own* (vendored)
 * Prism instance — but since the modules are already cached, their bodies are
 * no-ops and the grammars never get registered on the tokenizer. As a result
 * every `additionalLanguages` entry (bash, ini, hcl, ...) renders unhighlighted.
 *
 * To be robust against this, we always attach the grammars to the standalone
 * `prismjs` instance (whether cached or not) and then copy the resulting
 * language definitions onto prism-react-renderer's Prism object.
 */
export default function prismIncludeLanguages(PrismObject) {
  const {
    themeConfig: {prism},
  } = siteConfig
  const {additionalLanguages} = prism

  const PrismBefore = globalThis.Prism
  globalThis.Prism = Prism

  additionalLanguages.forEach((lang) => {
    if (lang === 'php') {
      // eslint-disable-next-line global-require
      require('prismjs/components/prism-markup-templating.js')
    }
    // eslint-disable-next-line global-require, import/no-dynamic-require
    require(`prismjs/components/prism-${lang}`)
  })

  additionalLanguages.forEach((lang) => {
    if (Prism.languages[lang]) {
      PrismObject.languages[lang] = Prism.languages[lang]
    }
  })

  delete globalThis.Prism
  if (typeof PrismBefore !== 'undefined') {
    globalThis.Prism = PrismBefore
  }
}
