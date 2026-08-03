import { Layout } from 'antd';
import { Suspense } from 'react';
import { Outlet } from 'react-router-dom';

import { layout } from '@/shared/design';
import { PageSkeleton } from '@/shared/ui';
import { Header } from '@/widgets/header';

export function AppLayout() {
  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header />

      <Layout.Content
        style={{
          maxWidth: layout.containerMaxWidth,
          width: '100%',
          margin: '0 auto',
          padding: 'var(--spacing-lg) var(--spacing-md)',
        }}
      >
        <Suspense fallback={<PageSkeleton />}>
          <Outlet />
        </Suspense>
      </Layout.Content>
    </Layout>
  );
}
