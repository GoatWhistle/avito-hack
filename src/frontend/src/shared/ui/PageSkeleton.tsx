import { Skeleton, Space } from 'antd';

const SKELETON_ROWS = 4;

export function PageSkeleton() {
  return (
    <Space direction="vertical" size="middle" style={{ width: '100%' }}>
      <Skeleton active title paragraph={{ rows: 1 }} />
      <Skeleton active title={false} paragraph={{ rows: SKELETON_ROWS }} />
    </Space>
  );
}
