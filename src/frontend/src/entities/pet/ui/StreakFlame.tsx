import { Tooltip } from 'antd';
import { useTranslation } from 'react-i18next';

import './pet-stats.css';

export interface StreakFlameProps {
  days: number;
  compact?: boolean;
}

interface FlameProps {
  size: number;
  active: boolean;
  label: string;
}

function Flame({ size, active, label }: FlameProps) {
  return (
    <svg
      className="pet-streak__flame"
      width={size}
      height={size * 1.15}
      viewBox="0 0 26 30"
      role="img"
      aria-label={label}
      focusable="false"
    >
      <path
        d="M 13 1 q 8 8 5 14 q 4 -1 4 -6 q 5 7 3 13 q -2 7 -12 7 q -10 0 -12 -7 q -2 -8 5 -14 q -1 5 2 6 q -3 -8 5 -13 z"
        fill={active ? 'var(--color-streak)' : 'var(--color-text-tertiary)'}
      />
      <path
        d="M 13 11 q 5 5 3 9 q -1 4 -3 4 q -2 0 -3 -4 q -2 -4 3 -9 z"
        fill={active ? 'var(--color-streak-core)' : 'var(--color-bg-track)'}
      />
    </svg>
  );
}

export function StreakFlame({ days, compact = false }: StreakFlameProps) {
  const { t } = useTranslation('pet');
  const isActive = days > 0;
  const tooltip = t('streak.tooltip', { count: days });

  if (compact) {
    return (
      <Tooltip title={tooltip}>
        <span
          className={isActive ? 'pet-streak--active' : 'pet-streak--idle'}
          style={{ display: 'inline-flex', alignItems: 'center', gap: 'var(--spacing-xs)' }}
        >
          <Flame size={18} active={isActive} label={t('streak.label')} />
          <span className="pet-stat__value">{days}</span>
        </span>
      </Tooltip>
    );
  }

  return (
    <section
      className={`pet-streak ${isActive ? 'pet-streak--active' : 'pet-streak--idle'}`}
      aria-label={tooltip}
    >
      <span className="pet-streak__icon">
        <Flame size={32} active={isActive} label={t('streak.label')} />
      </span>

      <span className="pet-streak__body">
        <span className="pet-streak__count">{days}</span>
        <span className="pet-streak__caption">
          {isActive ? t('streak.caption') : t('streak.captionIdle')}
        </span>
      </span>
    </section>
  );
}
