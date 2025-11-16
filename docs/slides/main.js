import Reveal from './node_modules/reveal.js/dist/reveal.esm.js';

import ExternalPlugin from './plugins/external.js';
import MarkdownPlugin from './node_modules/reveal.js/plugin/markdown/markdown.esm.js';
import HighlightPlugin from './node_modules/reveal.js/plugin/highlight/highlight.esm.js';
import SpeakerNotesPlugin from './node_modules/reveal.js/plugin/notes/notes.esm.js';
import MermaidPlugin from './plugins/mermaid.js';

Reveal.initialize({
  hash: true,
  plugins: [
    ExternalPlugin,
    MarkdownPlugin,
    HighlightPlugin,
    SpeakerNotesPlugin,
    MermaidPlugin,
  ],
  mermaid: {
    theme: 'dark',
    themeVariables: {
      darkMode: true,
    },
  },
});
