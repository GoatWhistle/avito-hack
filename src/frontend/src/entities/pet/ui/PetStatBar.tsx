import './pet-stats.css';

export type StatKind = 'satiety' | 'happiness';

export interface PetStatBarProps {
  kind: StatKind;
  label: string;
  value: number;
  low: boolean;
}

export function PetStatBar({ kind, label, value, low }: PetStatBarProps) {
  const clamped = Math.max(0, Math.min(100, value));
  const rounded = Math.round(clamped);

  return (
    <div className="pet-stat">
      <div className="pet-stat__head">
        <span className="pet-stat__label">{label}</span>
        <span className={low ? 'pet-stat__value pet-stat__value--low' : 'pet-stat__value'}>
          {rounded}
        </span>
      </div>

      <div
        className="pet-stat__track"
        role="meter"
        aria-label={label}
        aria-valuenow={rounded}
        aria-valuemin={0}
        aria-valuemax={100}
      >
        <div
          className={`pet-stat__fill pet-stat__fill--${low ? 'low' : kind}`}
          style={{ width: `${String(clamped)}%` }}
        />
      </div>
    </div>
  );
}
