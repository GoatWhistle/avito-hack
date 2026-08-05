import { useUnit } from 'effector-react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';

import { $pet, usePetProfile } from '@/entities/pet';
import { ROUTES } from '@/shared/config/routes';

function percentOf(xp: number, nextLevelXp: number): number {
  if (nextLevelXp <= 0) {
    return 100;
  }

  return Math.max(0, Math.min(100, Math.round((xp / nextLevelXp) * 100)));
}

export function LevelBadge() {
  const { t } = useTranslation('pet');
  const pet = useUnit($pet);
  const profile = usePetProfile(true);

  const level = pet?.level ?? profile.data?.level;
  const xp = pet?.xp ?? profile.data?.xp ?? 0;
  const nextLevelXp = pet?.nextLevelXp ?? profile.data?.xpToNextLevel ?? 0;

  if (level === undefined) {
    return null;
  }

  const percent = percentOf(xp, nextLevelXp);

  return (
    <Link
      className="header-level"
      to={ROUTES.pet}
      aria-label={t('xp.level', { level })}
      title={t('xp.progress', { xp, next: nextLevelXp })}
    >
      <span className="header-level__value">{level}</span>

      <span className="header-level__track" aria-hidden="true">
        <span className="header-level__fill" style={{ width: `${String(percent)}%` }} />
      </span>
    </Link>
  );
}
