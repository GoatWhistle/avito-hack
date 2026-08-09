import './pet-placeholder.css'

export interface PetPlaceholderProps {
  label: string
  sizeClass: string
  stage: string
  failed: boolean
  className?: string
}

export function PetPlaceholder({
  label,
  sizeClass,
  stage,
  failed,
  className,
}: PetPlaceholderProps) {
  const classes = [
    'pet-placeholder',
    sizeClass,
    failed ? 'pet-placeholder--failed' : '',
    className ?? '',
  ]
    .filter((value) => value !== '')
    .join(' ')

  return (
    <div
      className={classes}
      data-testid="pet-placeholder"
      data-stage={stage}
      data-failed={failed ? 'true' : 'false'}
      role="status"
      aria-live="polite"
      aria-label={label}
    >
      <span className="pet-placeholder__orb" aria-hidden="true" />
      <span className="pet-placeholder__caption">{label}</span>
    </div>
  )
}
