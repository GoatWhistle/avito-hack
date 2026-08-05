import { Layout } from 'antd';
import { Suspense } from 'react';
import { Outlet } from 'react-router-dom';

import { PageSkeleton } from '@/shared/ui';
import { Header } from '@/widgets/header';

import './app-layout.css';

export function AppLayout() {
  return (
    <Layout className="app-shell">
      <Header />

      <Layout.Content className="app-content" id="main">
        <Suspense fallback={<PageSkeleton />}>
          <Outlet />
        </Suspense>
      </Layout.Content>
    </Layout>
  );
}
