import { useTranslation } from 'react-i18next';

import './pet-stats.css';

export interface XpBarProps {
  level: number;
  xp: number;
  nextLevelXp: number;
  percent: number;
  remaining: number;
  maxLevel: boolean;
}

export function XpBar({ level, xp, nextLevelXp, percent, remaining, maxLevel }: XpBarProps) {
  const { t } = useTranslation('pet');
  const clamped = Math.max(0, Math.min(100, percent));

  return (
    <div className="pet-xp">
      <div className="pet-xp__head">
        <span className="pet-xp__level">{t('xp.level', { level })}</span>

        <span className="pet-xp__counter">
          {maxLevel ? t('xp.maxLevel') : t('xp.progress', { xp, next: nextLevelXp })}
        </span>
      </div>

      <div
        className="pet-xp__track"
        role="progressbar"
        aria-label={t('xp.label')}
        aria-valuenow={Math.round(clamped)}
        aria-valuemin={0}
        aria-valuemax={100}
      >
        <div className="pet-xp__fill" style={{ width: `${String(clamped)}%` }} />
      </div>

      {!maxLevel && (
        <p className="pet-xp__remaining">
          {t('xp.remaining', { count: remaining, level: level + 1 })}
        </p>
      )}
    </div>
  );
}
