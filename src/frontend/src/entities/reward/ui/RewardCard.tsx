import { CheckCircleFilled, GiftOutlined } from '@ant-design/icons';
import { Button, Card, Tag, Typography } from 'antd';
import { useTranslation } from 'react-i18next';

import type { RewardId, RewardProgress } from '../model/catalog';

import './reward-card.css';

export interface RewardCardProps {
  progress: RewardProgress;
  promocode: string | null;
  claiming: boolean;
  onClaim: (rewardId: RewardId) => void;
  highlight?: boolean;
}

const TITLE_KEY = {
  hatched_badge: 'catalog.hatched_badge.title',
  free_delivery: 'catalog.free_delivery.title',
  streak_freeze: 'catalog.streak_freeze.title',
  promo_discount: 'catalog.promo_discount.title',
  teen_skin: 'catalog.teen_skin.title',
  avito_scarf: 'catalog.avito_scarf.title',
  xp_booster: 'catalog.xp_booster.title',
  autoteka_discount: 'catalog.autoteka_discount.title',
  legend_skin: 'catalog.legend_skin.title',
} as const;

const DESCRIPTION_KEY = {
  hatched_badge: 'catalog.hatched_badge.description',
  free_delivery: 'catalog.free_delivery.description',
  streak_freeze: 'catalog.streak_freeze.description',
  promo_discount: 'catalog.promo_discount.description',
  teen_skin: 'catalog.teen_skin.description',
  avito_scarf: 'catalog.avito_scarf.description',
  xp_booster: 'catalog.xp_booster.description',
  autoteka_discount: 'catalog.autoteka_discount.description',
  legend_skin: 'catalog.legend_skin.description',
} as const;

const KIND_KEY = {
  promocode: 'kind.promocode',
  utility: 'kind.utility',
  cosmetic: 'kind.cosmetic',
} as const;

const CONDITION_KEY = {
  level: 'condition.level',
  streak: 'condition.streak',
} as const;

const REMAINING_KEY = {
  level: 'progress.remaining.level',
  streak: 'progress.remaining.streak',
} as const;

export function RewardCard({
  progress,
  promocode,
  claiming,
  onClaim,
  highlight = false,
}: RewardCardProps) {
  const { t } = useTranslation('reward');
  const { definition, unlocked, percent, current, remaining } = progress;

  return (
    <Card
      className={`reward-card${highlight ? ' reward-card--highlight' : ''}`}
      size="small"
      title={
        <span className="reward-card__title">
          {unlocked ? (
            <CheckCircleFilled className="reward-card__icon reward-card__icon--done" />
          ) : (
            <GiftOutlined className="reward-card__icon" />
          )}
          <span>{t(TITLE_KEY[definition.id])}</span>
        </span>
      }
      extra={<Tag>{t(KIND_KEY[definition.kind])}</Tag>}
    >
      <div className="reward-card__body">
        <p className="reward-card__description">{t(DESCRIPTION_KEY[definition.id])}</p>

        <p className="reward-card__condition">
          {t(CONDITION_KEY[definition.condition], { target: definition.target })}
        </p>

        <div
          className="reward-card__track"
          role="progressbar"
          aria-label={t('progress.label')}
          aria-valuenow={percent}
          aria-valuemin={0}
          aria-valuemax={100}
        >
          <div
            className={`reward-card__fill${unlocked ? ' reward-card__fill--done' : ''}`}
            style={{ width: `${String(percent)}%` }}
          />
        </div>

        <p className={`reward-card__status${unlocked ? ' reward-card__status--done' : ''}`}>
          {unlocked
            ? t('progress.unlocked')
            : t(REMAINING_KEY[definition.condition], {
                count: remaining,
                current,
                target: definition.target,
              })}
        </p>

        {promocode === null ? (
          <Button
            type={unlocked ? 'primary' : 'default'}
            block
            disabled={!unlocked}
            loading={claiming}
            onClick={() => {
              onClaim(definition.id);
            }}
          >
            {t('actions.claim')}
          </Button>
        ) : (
          <div className="reward-card__code">
            <span className="reward-card__code-label">{t('promocode.label')}</span>
            <Typography.Text className="reward-card__code-value" strong copyable>
              {promocode}
            </Typography.Text>
            <span className="reward-card__code-hint">{t('promocode.hint')}</span>
          </div>
        )}
      </div>
    </Card>
  );
}
