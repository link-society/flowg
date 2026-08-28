import OpenTelemetryImage from '@site/src/assets/opentelemetry.svg'

import styles from './styles.module.css'


export default [
  {
    key: 'http',
    element: (
      <>
        <span className={styles.cardIcon} aria-hidden="true">
          HTTP
        </span>
        <span className={styles.cardLabel}>
          Webhooks &amp; APIs
        </span>
      </>
    ),
  },
  {
    key: 'syslog',
    element: (
      <>
        <span className={styles.cardIcon} aria-hidden="true">
          &lt;/&gt;
        </span>
        <span className={styles.cardLabel}>
          Syslog (TCP &amp; UDP)
        </span>
      </>
    ),
  },
  {
    key: 'file',
    element: (
      <>
        <span className={styles.cardIcon} aria-hidden="true">
          TXT
        </span>
        <span className={styles.cardLabel}>
          Log files
        </span>
      </>
    ),
  },
  {
    key: 'otlp',
    element: (
      <OpenTelemetryImage title="OpenTelemetry" />
    ),
  },
]
