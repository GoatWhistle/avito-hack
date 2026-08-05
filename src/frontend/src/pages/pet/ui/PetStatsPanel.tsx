import { ArrowRightOutlined } from '@ant-design/icons';
import { Button, Card } from 'antd';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';

import { PetStatBar, XpBar } from '@/entities/pet';
import { ROUTES } from '@/shared/config/routes';

const SATIETY_LOW = 30;
const HAPPINESS_LOW = 30;

export interface PetStatsPanelProps {
  level: number;
  xp: number;
  nextLevelXp: number;
  xpPercent: number;
  xpLeft: number;
  maxLevel: boolean;
  satiety: number;
  happiness: number;
}

export function PetStatsPanel({
  level,
  xp,
  nextLevelXp,
  xpPercent,
  xpLeft,
  maxLevel,
  satiety,
  happiness,
}: PetStatsPanelProps) {
  const { t } = useTranslation('pet');

  return (
    <>
      <Card>
        <div className="pet-vitals">
          <XpBar
            level={level}
            xp={xp}
            nextLevelXp={nextLevelXp}
            percent={xpPercent}
            remaining={xpLeft}
            maxLevel={maxLevel}
          />

          <PetStatBar
            kind="satiety"
            label={t('stats.satiety')}
            value={satiety}
            low={satiety < SATIETY_LOW}
          />

          <PetStatBar
            kind="happiness"
            label={t('stats.happiness')}
            value={happiness}
            low={happiness < HAPPINESS_LOW}
          />
        </div>
      </Card>

      <Card title={<span className="app-section-title">{t('nextSteps.title')}</span>}>
        <div className="pet-next">
          <div className="pet-next__row">
            <p className="pet-next__text">{t('nextSteps.publish')}</p>
            <Link to={ROUTES.itemCreate}>
              <Button type="primary" icon={<ArrowRightOutlined />} iconPosition="end">
                {t('nextSteps.publishAction')}
              </Button>
            </Link>
          </div>

          <div className="pet-next__row">
            <p className="pet-next__text">{t('nextSteps.rewards')}</p>
            <Link to={ROUTES.rewards}>
              <Button icon={<ArrowRightOutlined />} iconPosition="end">
                {t('nextSteps.rewardsAction')}
              </Button>
            </Link>
          </div>
        </div>
      </Card>
    </>
  );
}
