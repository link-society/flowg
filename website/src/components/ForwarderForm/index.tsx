import React from 'react'

import styles from './styles.module.css'

// keyboard icon = static mode; code icon = expression mode
const STATIC_ICON =
  'M21 3.01H3c-1.1 0-2 .9-2 2V9h2V4.99h18v14.03H3V15H1v4.01c0 1.1.9 1.98 2 1.98h18c1.1 0 2-.88 2-1.98v-14c0-1.11-.9-2-2-2M11 16l4-4-4-4v3H1v2h10z'
const EXPR_ICON =
  'M9.4 16.6 4.8 12l4.6-4.6L8 6l-6 6 6 6zm5.2 0 4.6-4.6-4.6-4.6L16 6l6 6-6 6z'
const FORWARDER_ICON =
  'M20 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h9v-2H4V8l8 5 8-5v5h2V6c0-1.1-.9-2-2-2m-8 7L4 6h16zm7 4 4 4-4 4v-3h-4v-2h4z'

export type ForwarderFormField =
  | { type?: 'text' | 'number'; label: string; value: string | number; tooltip?: string }
  | { type: 'password'; label: string; value?: string; tooltip?: string }
  | { type: 'textarea'; label: string; value: string; tooltip?: string }
  | {
      /** Dynamic field: prefix value with `@expr:` for expression mode. */
      type: 'dynamic'
      label: string
      value: string
      tooltip?: string
    }
  | {
      /** Key-value pairs (mirrors InputKeyValue from the web app). */
      type: 'keyvalue'
      label?: string
      keyLabel?: string
      valueLabel?: string
      pairs: [string, string][]
      tooltip?: string
    }
  | { type: 'checkbox'; label: string; checked: boolean; tooltip?: string }
  | { type: 'divider' }

export type ForwarderFormProps = {
  /** Displayed in the disabled "Forwarder Type" field at the top. */
  forwarderType: string
  fields: ForwarderFormField[]
}

function Tooltip({ text }: { text: string }) {
  return (
    <span className={styles.tooltip}>
      <span className={styles.tooltipIcon}>?</span>
      <span className={styles.tooltipPopup}>{text}</span>
    </span>
  )
}

function Field({ spec }: { spec: ForwarderFormField }) {
  if (spec.type === 'divider') {
    return <hr className={styles.divider} />
  }

  if (spec.type === 'dynamic') {
    const isExpr = spec.value.startsWith('@expr:')
    const icon = isExpr ? EXPR_ICON : STATIC_ICON
    const displayValue = isExpr ? spec.value.slice(6) : spec.value

    return (
      <div className={styles.fieldRow}>
        <div className={`${styles.field} ${styles.dynamicField}`}>
          <span className={styles.label}>{spec.label}</span>
          <span className={styles.dynamicMode}>
            <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor">
              <path d={icon} />
            </svg>
          </span>
          <span className={styles.dynamicSeparator} />
          <span className={styles.dynamicValue}>{displayValue}</span>
        </div>
        {spec.tooltip && <Tooltip text={spec.tooltip} />}
      </div>
    )
  }

  if (spec.type === 'keyvalue') {
    const keyLabel = spec.keyLabel ?? 'Key'
    const valueLabel = spec.valueLabel ?? 'Value'

    return (
      <>
        {spec.pairs.map(([key, value], i) => (
          <div key={i} className={styles.kvRow}>
            <div className={styles.kvCell}>
              <span className={styles.kvCellLabel}>{keyLabel}</span>
              <span className={styles.kvCellValue}>{key}</span>
            </div>
            <div className={styles.kvCell}>
              <span className={styles.kvCellLabel}>{valueLabel}</span>
              <span className={styles.kvCellValue}>{value}</span>
            </div>
            <span className={`${styles.kvBtn} ${styles.kvBtnDelete}`}>
              <svg viewBox="0 0 24 24" width="18" height="18" fill="currentColor">
                <path d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6zM19 4h-3.5l-1-1h-5l-1 1H5v2h14z" />
              </svg>
            </span>
          </div>
        ))}
        <div className={styles.kvRow}>
          <div className={styles.kvCell}>
            <span className={styles.kvCellLabel}>{keyLabel}</span>
          </div>
          <div className={styles.kvCell}>
            <span className={styles.kvCellLabel}>{valueLabel}</span>
          </div>
          <span className={`${styles.kvBtn} ${styles.kvBtnAdd}`}>
            <svg viewBox="0 0 24 24" width="18" height="18" fill="currentColor">
              <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z" />
            </svg>
          </span>
          {spec.tooltip && <Tooltip text={spec.tooltip} />}
        </div>
      </>
    )
  }

  if (spec.type === 'checkbox') {
    return (
      <span className={styles.checkboxField}>
        <span
          className={`${styles.checkbox} ${spec.checked ? styles.checkboxChecked : ''}`}
        >
          {spec.checked && (
            <svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor">
              <path d="M9 16.17 4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z" />
            </svg>
          )}
        </span>
        {spec.label}
        {spec.tooltip && <Tooltip text={spec.tooltip} />}
      </span>
    )
  }

  const displayValue =
    spec.type === 'password' ? '\u2022'.repeat(12) : String(spec.value ?? '')

  return (
    <div className={styles.fieldRow}>
      <div className={styles.field}>
        <span className={styles.label}>{spec.label}</span>
        <span
          className={`${styles.value} ${spec.type === 'textarea' ? styles.valueTextarea : ''}`}
        >
          {displayValue}
        </span>
      </div>
      {spec.tooltip && <Tooltip text={spec.tooltip} />}
    </div>
  )
}

export default function ForwarderForm({
  forwarderType,
  fields,
}: ForwarderFormProps) {
  return (
    <div className={styles.form}>
      <div className={`${styles.field} ${styles.fieldDisabled}`}>
        <span className={styles.label}>Forwarder Type</span>
        <span className={`${styles.value} ${styles.valueDisabled}`}>
          <svg
            className={styles.typeIcon}
            viewBox="0 0 24 24"
            width="20"
            height="20"
            fill="currentColor"
          >
            <path d={FORWARDER_ICON} />
          </svg>
          {forwarderType}
        </span>
      </div>

      {fields.map((spec, i) => (
        <Field key={i} spec={spec} />
      ))}
    </div>
  )
}
