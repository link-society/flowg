export default {
  plain: {
    backgroundColor: "#272822",
    color: "#F8F8F2"
  },
  styles: [
    {
      types: ["comment", "prolog", "doctype", "cdata"],
      style: {
        color: "#778090"
      }
    },
    {
      types: ["punctuation"],
      style: {
        color: "#F8F8F2"
      }
    },
    {
      types: ["namespace"],
      style: {
        opacity: 0.7
      }
    },
    {
      types: ["property", "tag", "constant", "symbol", "deleted"],
      style: {
        color: "#F92672"
      }
    },
    {
      types: ["boolean", "number"],
      style: {
        color: "#AE81FF"
      }
    },
    {
      types: ["selector", "attr-name", "string", "char", "builtin", "inserted"],
      style: {
        color: "#A6E22E"
      }
    },
    {
      types: ["operator", "entity", "url", "variable"],
      style: {
        color: "#F8F8F2"
      }
    },
    {
      types: ["atrule", "attr-value", "function"],
      style: {
        color: "#E6DB74"
      }
    },
    {
      types: ["keyword", "key"],
      style: {
        color: "#F92672"
      }
    },
    {
      types: ["regex", "important"],
      style: {
        color: "#FD971F"
      }
    },
    {
      types: ["important", "bold"],
      style: {
        fontWeight: "bold"
      }
    },
    {
      types: ["italic"],
      style: {
        fontStyle: "italic"
      }
    },
    {
      types: ["entity"],
      style: {
        cursor: "help"
      }
    },
    {
      types: ["string"],
      languages: ["css"],
      style: {
        color: "#F8F8F2"
      }
    },
    {
      types: ["attr-name"],
      languages: ["ini"],
      style: {
        color: "#F92672"
      }
    },
    {
      types: ["variable"],
      languages: ["bash"],
      style: {
        color: "#F92672"
      }
    }
  ]
}