const editorId = arguments[0]
const code = arguments[1]

const getWrapper = (ed) => ed.getContainerDomNode().parentElement
const hasId = (ed) => getWrapper(ed).id === editorId

const [editor] = window.monaco.editor.getEditors().filter(hasId)
editor.getModel().setValue(code)
