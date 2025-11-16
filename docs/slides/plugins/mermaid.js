/*!
 * reveal.js Mermaid plugin
 */

import mermaid from "../node_modules/mermaid/dist/mermaid.esm.mjs";
import { icons } from "../node_modules/@iconify-json/logos/index.mjs";

async function renderMermaid({ el, beforeRender, afterRender }) {
  const beforeRenderRes = await beforeRender?.(el);

  if (beforeRenderRes === false) {
    return;
  }

  // Using textContent not innerHTML, because innerHTML will get escaped code (eg: get --&gt; instead of -->).
  const graphDefinition = el.textContent.trim();

  try {
    const { svg: svgCode } = await mermaid.render(
      `mermaid-${Math.random().toString(36).substring(2)}`,
      graphDefinition
    );
    el.innerHTML = svgCode;

    await afterRender?.(el);
  } catch (error) {
    let errorStr = "";
    if (error?.str) {
      // From mermaid 9.1.4, error.message does not exists anymore
      errorStr = error.str;
    }
    if (error?.message) {
      errorStr = error.message;
    }
    console.error(errorStr, { error, graphDefinition, el });
    el.innerHTML = errorStr;
  }
}

function getRenderMermaidEl({ beforeRender, afterRender }) {
  return function renderMermaidEl(el) {
    return renderMermaid({
      el,
      beforeRender,
      afterRender,
    });
  };
}

const Plugin = {
  id: "mermaid",

  init: function (reveal) {
    const { ...mermaidConfig } = reveal.getConfig().mermaid || {};
    const { ...mermaidPluginConfig } = reveal.getConfig().mermaidPlugin || {};

    const renderMermaidEl = getRenderMermaidEl({
      beforeRender: mermaidPluginConfig.beforeRender,
      afterRender: mermaidPluginConfig.afterRender,
    });

    mermaid.initialize({
      // The node size will be calculated incorrectly if set `startOnLoad: start`,
      // so we need to manually render.
      startOnLoad: false,
      ...mermaidConfig,
    });

    mermaid.registerIconPacks([
      {
        name: icons.prefix,
        icons,
      },
    ]);

    const mermaidEls = reveal.getRevealElement().querySelectorAll(".mermaid");

    Array.from(mermaidEls).forEach(function (el) {
      renderMermaidEl(el);
    });
  },
};

export default () => Plugin;
