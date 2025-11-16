/*
* external.js by Jan Schoepke <janschoepke@me.com> - MIT license
* modernized for Reveal.js 5 by Stephan Heunis - MIT license
*/

// Plugin structure inspired by reveal markdown plugin
const Plugin = () => {
    let deck;
    let options = {
        // async option is moot after migrating to fetch api
        mapAttributes: []
    };

    // getTarget function unchanged
    const getTarget = function(node) {
        let url = node.getAttribute('data-external') || '';
        let isReplace = false;
        if (url === '') {
            url = node.getAttribute('data-external-replace') || '';
            isReplace = true;
        }
        if (url.length > 0) {
            const r = url.match(/^([^#]+)(?:#(.+))?/);
            return {
                url: r[1] || '',
                fragment: r[2] || '',
                isReplace: isReplace
            };
        }
        return null;
    };

    // convertUrl function unchanged
    const convertUrl = function(src, path) {
        if (path !== '' && src.indexOf('.') === 0) {
            return path + '/' + src;
        }
        return src;
    };

    // convertAttributes function unchanged
    const convertAttributes = function(attributeName, container, path) {
        const nodes = container.querySelectorAll('[' + attributeName + ']');
        if (container.getAttribute(attributeName)) {
            container.setAttribute(
                attributeName,
                convertUrl(container.getAttribute(attributeName), path)
            );
        }
        for (let i = 0; i < nodes.length; i++) {
            nodes[i].setAttribute(
                attributeName,
                convertUrl(nodes[i].getAttribute(attributeName), path)
            );
        }
    };

    // convertUrls function unchanged
    const convertUrls = function(container, path) {
        for (let i = 0; i < options.mapAttributes.length; i++) {
            convertAttributes(options.mapAttributes[i], container, path);
        }
    };

    // updateSection function changed:
    // - async function
    // - migrated to fetch api
    // - await recursive slides
    const updateSection = async function(section, target, path) {
        const url = path !== '' ? path + '/' + target.url : target.url;

        try {
            const response = await fetch(url);
            if (!response.ok) {
                throw new Error('HTTP error ' + response.status);
            }

            const text = await response.text();
            const basePath = url.substr(0, url.lastIndexOf('/'));
            const html = new DOMParser().parseFromString(text, 'text/html');
            // Account for possibility of body not existing
            const source = html.querySelector('body') || html;
            let nodes = target.fragment !== ''
                ? html.querySelectorAll(target.fragment)
                : source.childNodes;
			// ensure array with only elements, no stray text nodes
			nodes = Array.from(nodes).filter(n => n.nodeType === 1);
            if (!target.isReplace) {
                section.innerHTML = '';
            }
            for (const node of nodes) {
                convertUrls(node, basePath);
                const imported = document.importNode(node, true);
                target.isReplace
                    ? section.parentNode.insertBefore(imported, section)
                    : section.appendChild(imported);

                if (imported instanceof Element) {
                    // For recursive loading, we need to await otherwise recursive
                    // slides might only load after reveal syncs slides again
                    await loadExternal(imported, basePath);
                }
            }
            if (target.isReplace) {
                section.parentNode.removeChild(section);
            }
        } catch (error) {
            console.log('Failed to fetch ' + url + ': ' + error.message);
        }
    };

    // loadExternal function changed:
    // - async function
    // - migrated to fetch api
    // - await recursive slides
    // - build up array of promises and await all
    async function loadExternal(container, path) {
        let promises = [];
        path = typeof path === 'undefined' ? '' : path;
        if (
            container instanceof Element &&
            (container.getAttribute('data-external') ||
                container.getAttribute('data-external-replace'))
        ) {
            const target = getTarget(container);
            if (target) {
                promises.push(updateSection(container, target, path));
            }
        } else {
            const sections = container.querySelectorAll(
                '[data-external], [data-external-replace]'
            );
            for (let i = 0; i < sections.length; i++) {
                const section = sections[i];
                const target = getTarget(section);
                if (target) {
                    promises.push(updateSection(section, target, path));
                }
            }
        }
        await Promise.all(promises);
    }

    // Plugin returns:
    return {
        id: 'external',
        init: function(reveal) {
            deck = reveal;
            const cfg = deck.getConfig().external || {};
            options = {
                mapAttributes: Array.isArray(cfg.mapAttributes)
                ? cfg.mapAttributes
                : (cfg.mapAttributes ? ['src'] : [])
            };
            // return promise, which latest reveal can handle
            return loadExternal(deck.getRevealElement()).then(() => {
                // recalculates deck layout after loading external slides
                deck.layout();
            });
        }
    };
};

export default Plugin;
