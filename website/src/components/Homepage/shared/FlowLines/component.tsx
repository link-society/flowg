import clsx from 'clsx'

import type { FlowLinesProps } from './types'

import styles from './styles.module.css'

const FlowLines = ({ className }: FlowLinesProps) => {
  return (
    <svg
      className={clsx(styles.flowArt, className)}
      viewBox="0 0 1440 580"
      fill="none"
      aria-hidden="true"
      focusable="false"
    >
      <g stroke="currentColor" strokeWidth="1.2">
        <path d="M0 95h125q40 0 40 40v175q0 40 40 40h70" />
        <path d="M0 165h75q40 0 40 40v210q0 40 40 40h135" />
        <path d="M0 340h35q30 0 30 30v130q0 30 30 30h220" />
        <path d="M1440 100h-85q-40 0-40 40v115q0 40-40 40h-70" />
        <path d="M1440 200h-35q-40 0-40 40v140q0 40-40 40h-110" />
        <path d="M1440 330h-10q-20 0-20 20v150q0 30-30 30h-190" />
      </g>
      <g
        fill="var(--fg-art-node-fill, white)"
        stroke="currentColor"
        strokeWidth="2"
      >
        <circle cx="165" cy="200" r="5" />
        <circle cx="115" cy="365" r="5" />
        <circle cx="65" cy="435" r="5" />
        <circle cx="1315" cy="200" r="5" />
        <circle cx="1365" cy="340" r="5" />
        <circle cx="1410" cy="435" r="5" />
      </g>
    </svg>
  )
}

export default FlowLines
