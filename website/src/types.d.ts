// Docusaurus only declares *.svg / *.css / *.md modules; add *.png for image imports.
declare module '*.png' {
  const src: string
  export default src
}
