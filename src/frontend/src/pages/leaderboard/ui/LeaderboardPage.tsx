import { Table, Tag } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { useUnit } from 'effector-react';
import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';

import { useLeaderboard, type LeaderboardEntry } from '@/entities/leaderboard';
import { StreakFlame } from '@/entities/pet';
import { $isAuthenticated, $user } from '@/entities/session';
import { EmptyState, ErrorState, PageSkeleton } from '@/shared/ui';

import './leaderboard-page.css';

const MEDALS = ['leaderboard-rank--gold', 'leaderboard-rank--silver', 'leaderboard-rank--bronze'];

function rankClass(rank: number): string {
  return rank <= MEDALS.length ? `leaderboard-rank ${MEDALS[rank - 1]}` : 'leaderboard-rank';
}

export function LeaderboardPage() {
  const { t } = useTranslation('leaderboard');
  const isAuthenticated = useUnit($isAuthenticated);
  const user = useUnit($user);
  const query = useLeaderboard(user?.id ?? null, isAuthenticated);

  const columns = useMemo<ColumnsType<LeaderboardEntry>>(
    () => [
      {
        title: t('columns.rank'),
        dataIndex: 'rank',
        key: 'rank',
        width: 72,
        render: (rank: number) => <span className={rankClass(rank)}>{rank}</span>,
      },
      {
        title: t('columns.name'),
        dataIndex: 'name',
        key: 'name',
        render: (name: string, entry) => (
          <span className="leaderboard-name">
            <span className="leaderboard-name__text">{name}</span>
            {entry.isMe && <Tag color="blue">{t('you')}</Tag>}
          </span>
        ),
      },
      {
        title: t('columns.level'),
        dataIndex: 'level',
        key: 'level',
        width: 96,
        responsive: ['sm'],
        render: (level: number) => <span className="leaderboard-num">{level}</span>,
      },
      {
        title: t('columns.xp'),
        dataIndex: 'xp',
        key: 'xp',
        width: 116,
        render: (xp: number) => <span className="leaderboard-num">{xp}</span>,
      },
      {
        title: t('columns.streak'),
        dataIndex: 'streakDays',
        key: 'streakDays',
        width: 104,
        responsive: ['md'],
        render: (days: number) => <StreakFlame days={days} compact />,
      },
    ],
    [t],
  );

  if (!isAuthenticated) {
    return <EmptyState description={t('guard.signIn')} />;
  }

  if (query.isError) {
    return (
      <ErrorState
        error={query.error}
        onRetry={() => {
          void query.refetch();
        }}
      />
    );
  }

  if (query.isPending) {
    return <PageSkeleton />;
  }

  const entries = query.data.entries;
  const myRank = query.data.myRank;

  return (
    <div className="app-stack">
      <header className="leaderboard-head">
        <h1 className="app-page-title">{t('title')}</h1>
        {myRank != null && (
          <p className="leaderboard-head__rank">
            {t('myRankLabel')} <strong>#{myRank}</strong>
          </p>
        )}
      </header>

      {entries.length === 0 ? (
        <EmptyState description={t('empty')} />
      ) : (
        <div className="leaderboard-table">
          <Table
            rowKey="userId"
            columns={columns}
            dataSource={entries}
            pagination={false}
            size="middle"
            rowClassName={(entry) => (entry.isMe ? 'leaderboard-row-me' : '')}
          />
        </div>
      )}
    </div>
  );
}
