import { useUnit } from 'effector-react';
import { useTranslation } from 'react-i18next';

import {
  PetAvatar,
  StreakFlame,
  usePetProfile,
  type PetState,
  type RaccoonProfile,
} from '@/entities/pet';
import { $isAuthenticated } from '@/entities/session';
import { PetActions, usePetLive } from '@/features/pet-actions';
import { EmptyState, PageSkeleton } from '@/shared/ui';

import { ConnectionBadge } from './ConnectionBadge';
import './pet-page.css';
import { PetStatsPanel } from './PetStatsPanel';

interface Summary {
  name: string | null;
  level: number;
  xp: number;
  nextLevelXp: number;
  streak: number;
  satiety: number;
  happiness: number;
}

function summarize(pet: PetState | null, profile: RaccoonProfile | undefined): Summary {
  if (pet !== null) {
    return {
      name: pet.name === '' ? null : pet.name,
      level: pet.level,
      xp: pet.xp,
      nextLevelXp: pet.nextLevelXp,
      streak: pet.streakDays,
      satiety: pet.satiety,
      happiness: pet.happiness,
    };
  }

  return {
    name: profile?.name ?? null,
    level: profile?.level ?? 1,
    xp: profile?.xp ?? 0,
    nextLevelXp: profile?.xpToNextLevel ?? 0,
    streak: profile?.currentStreak ?? 0,
    satiety: 0,
    happiness: 0,
  };
}

export function PetPage() {
  const { t } = useTranslation('pet');
  const isAuthenticated = useUnit($isAuthenticated);
  const live = usePetLive(isAuthenticated);
  const profile = usePetProfile(isAuthenticated);

  if (!isAuthenticated) {
    return <EmptyState description={t('guard.signIn')} />;
  }

  if (live.pet === null && profile.isPending) {
    return <PageSkeleton />;
  }

  const summary = summarize(live.pet, profile.data);

  return (
    <div className="pet-page">
      <section className="pet-stage" aria-labelledby="pet-name">
        <div className="pet-stage__top">
          <h1 className="pet-stage__name" id="pet-name">
            {summary.name ?? t('title')}
          </h1>

          <ConnectionBadge status={live.status} />
        </div>

        <div className="pet-stage__avatar">
          <PetAvatar
            stage={live.pet?.stage ?? 'egg'}
            emotion={live.emotion}
            satiety={summary.satiety}
            happiness={summary.happiness}
            onStroke={live.stroke}
          />
        </div>

        <p className="pet-stage__hint" aria-live="polite">
          {t(`hint.${live.emotion}`)}
        </p>

        <PetActions status={live.status} onStroke={live.stroke} />
      </section>

      <aside className="pet-aside">
        <StreakFlame days={summary.streak} />

        <PetStatsPanel
          level={summary.level}
          xp={summary.xp}
          nextLevelXp={summary.nextLevelXp}
          xpPercent={live.xpPercent}
          xpLeft={live.xpLeft}
          maxLevel={live.maxLevel}
          satiety={summary.satiety}
          happiness={summary.happiness}
        />
      </aside>
    </div>
  );
}
