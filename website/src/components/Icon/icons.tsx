const paths = {
  arrow: <path d="M5 12h14m-6-6 6 6-6 6" />,
  code: <path d="m8 7-5 5 5 5m8-10 5 5-5 5m-3-14-2 18" />,
  file: (
    <>
      <path d="M14 3H5v18h14V8Z" />
      <path d="M14 3v5h5M8 12h8m-8 4h6" />
    </>
  ),
  filter: <path d="M3 5h18l-7 8v6l-4 2v-8Z" />,
  route: (
    <path d="M6 5v10a3 3 0 0 0 3 3h9M6 11h8a4 4 0 0 0 4-4V5m-3 3 3-3 3 3m-6 7 3 3-3 3" />
  ),
  cloud: (
    <path d="M7 18a5 5 0 1 1 .4-10A7 7 0 0 1 21 11a3.5 3.5 0 0 1-1.5 7Z" />
  ),
  search: (
    <>
      <circle cx="10.5" cy="10.5" r="6.5" />
      <path d="m16 16 5 5" />
    </>
  ),
  spark: (
    <>
      <path d="m12 3 2.5 6.5L21 12l-6.5 2.5L12 21l-2.5-6.5L3 12l6.5-2.5Z" />
      <path d="M20 2v4m-2-2h4" />
    </>
  ),
  key: (
    <>
      <circle cx="8" cy="9" r="5" />
      <path d="m12 13 8 8m-3-3 3-3m-6 0 3-3" />
    </>
  ),
  upgrade: (
    <>
      <path d="M20 8a8 8 0 0 0-14-3L3 8m0-5v5h5M4 16a8 8 0 0 0 14 3l3-3m0 5v-5h-5" />
    </>
  ),
  maintenance: (
    <path d="M14 6a5 5 0 0 0-6 6l-5 5a3 3 0 0 0 4 4l5-5a5 5 0 0 0 6-6l-3 3-4-4Z" />
  ),
  support: (
    <>
      <path d="M4 13v-2a8 8 0 0 1 16 0v6a4 4 0 0 1-4 4h-3" />
      <rect x="2" y="11" width="5" height="7" rx="2" />
      <rect x="17" y="11" width="5" height="7" rx="2" />
    </>
  ),
  check: <path d="m5 12 4 4L19 6" />,
  terminal: (
    <>
      <rect x="3" y="4" width="18" height="16" rx="3" />
      <path d="m7 9 3 3-3 3m6 0h4" />
    </>
  ),
}

export default paths
