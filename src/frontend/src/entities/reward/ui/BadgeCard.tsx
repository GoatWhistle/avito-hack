import { TrophyFilled } from '@ant-design/icons';
import { Card, Space, Typography } from 'antd';
import { useTranslation } from 'react-i18next';

import { petPalette } from '@/shared/design';
import { formatDate } from '@/shared/lib/format';

export interface BadgeCardProps {
  name: string;
  description: string;
  iconUrl: string;
  earnedAt: string | null;
}

const ICON_SIZE = 40;

export function BadgeCard({ name, description, iconUrl, earnedAt }: BadgeCardProps) {
  const { t, i18n } = useTranslation('reward');

  return (
    <Card size="small">
      <Space align="start" size={12}>
        {iconUrl === '' ? (
          <TrophyFilled style={{ fontSize: ICON_SIZE, color: petPalette.streakCore }} />
        ) : (
          <img src={iconUrl} alt={name} width={ICON_SIZE} height={ICON_SIZE} />
        )}

        <Space direction="vertical" size={2}>
          <Typography.Text strong>{name}</Typography.Text>
          <Typography.Text type="secondary">{description}</Typography.Text>
          <Typography.Text type="secondary">
            {earnedAt === null
              ? t('badge.locked')
              : t('badge.earnedAt', { date: formatDate(earnedAt, i18n.language) })}
          </Typography.Text>
        </Space>
      </Space>
    </Card>
  );
}
