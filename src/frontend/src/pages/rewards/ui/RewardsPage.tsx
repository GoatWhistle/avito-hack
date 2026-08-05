import { useUnit } from 'effector-react';
import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';

import { useBadges, usePetProfile, $pet } from '@/entities/pet';
import { buildRewardProgress, BadgeCard, RewardCard, type RewardProgress } from '@/entities/reward';
import { $isAuthenticated } from '@/entities/session';
import { useRewardClaim } from '@/features/reward-claim';
import { EmptyState, ErrorState, PageSkeleton } from '@/shared/ui';

import './rewards-page.css';

interface Groups {
  unlocked: RewardProgress[];
  close: RewardProgress[];
  far: RewardProgress[];
}

const CLOSE_PERCENT = 50;

const GROUP_KEY = {
  close: 'group.close',
  unlocked: 'group.unlocked',
  far: 'group.far',
} as const;

type GroupKey = keyof typeof GROUP_KEY;

function group(items: RewardProgress[]): Groups {
  const unlocked = items.filter((item) => item.unlocked);
  const locked = items
    .filter((item) => !item.unlocked)
    .sort((a, b) => b.percent - a.percent || a.remaining - b.remaining);

  return {
    unlocked,
    close: locked.filter((item) => item.percent >= CLOSE_PERCENT),
    far: locked.filter((item) => item.percent < CLOSE_PERCENT),
  };
}

export function RewardsPage() {
  const { t } = useTranslation('reward');
  const isAuthenticated = useUnit($isAuthenticated);
  const pet = useUnit($pet);
  const profile = usePetProfile(isAuthenticated);
  const badges = useBadges(isAuthenticated);
  const { promocodes, pendingId, claim } = useRewardClaim();

  const level = pet?.level ?? profile.data?.level ?? 1;
  const streak = pet?.streakDays ?? profile.data?.currentStreak ?? 0;
  const groups = useMemo(() => group(buildRewardProgress(level, streak)), [level, streak]);

  if (!isAuthenticated) {
    return <EmptyState description={t('guard.signIn')} />;
  }

  if (profile.isError) {
    return (
      <ErrorState
        error={profile.error}
        onRetry={() => {
          void profile.refetch();
        }}
      />
    );
  }

  if (profile.isPending) {
    return <PageSkeleton />;
  }

  const earnedBadges = badges.data ?? profile.data.badges;

  const renderGroup = (key: GroupKey, items: RewardProgress[], highlight: boolean) =>
    items.length > 0 && (
      <section className="rewards-group" aria-labelledby={`rewards-${key}`}>
        <h2 className="app-section-title" id={`rewards-${key}`}>
          {t(GROUP_KEY[key])}
        </h2>

        <div className="rewards-grid">
          {items.map((item) => (
            <RewardCard
              key={item.definition.id}
              progress={item}
              promocode={promocodes[item.definition.id] ?? null}
              claiming={pendingId === item.definition.id}
              onClaim={claim}
              highlight={highlight}
            />
          ))}
        </div>
      </section>
    );

  return (
    <div className="app-stack">
      <header>
        <h1 className="app-page-title">{t('title')}</h1>
        <p className="rewards-subtitle">{t('subtitle')}</p>
      </header>

      {renderGroup('close', groups.close, true)}
      {renderGroup('unlocked', groups.unlocked, false)}
      {renderGroup('far', groups.far, false)}

      <section className="rewards-group" aria-labelledby="rewards-badges">
        <h2 className="app-section-title" id="rewards-badges">
          {t('badge.title')}
        </h2>

        {earnedBadges.length === 0 ? (
          <EmptyState description={t('badge.empty')} />
        ) : (
          <div className="rewards-grid">
            {earnedBadges.map((badge) => (
              <BadgeCard
                key={badge.id}
                name={badge.name}
                description={badge.description}
                iconUrl={badge.iconUrl}
                earnedAt={badge.earnedAt}
              />
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
