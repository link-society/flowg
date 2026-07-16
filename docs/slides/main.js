import Reveal from './node_modules/reveal.js/dist/reveal.mjs';

import ExternalPlugin from './plugins/external.js';
import MarkdownPlugin from './node_modules/reveal.js/dist/plugin/markdown.mjs';
import HighlightPlugin from './node_modules/reveal.js/dist/plugin/highlight.mjs';
import SpeakerNotesPlugin from './node_modules/reveal.js/dist/plugin/notes.mjs';
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
