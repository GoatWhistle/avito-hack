import { Skeleton } from 'antd';

import './shared-ui.css';

const CARD_COUNT = [0, 1, 2];

export interface PageSkeletonProps {
  cards?: boolean;
}

export function PageSkeleton({ cards = true }: PageSkeletonProps) {
  return (
    <div className="page-skeleton" aria-busy="true">
      <Skeleton active title paragraph={{ rows: 1 }} />

      {cards && (
        <div className="page-skeleton__grid">
          {CARD_COUNT.map((index) => (
            <div className="page-skeleton__card" key={index}>
              <Skeleton active title paragraph={{ rows: 3 }} />
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
