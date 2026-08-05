import { lazy } from 'react';
import { createBrowserRouter, Navigate } from 'react-router-dom';

import { ROUTES, ROUTE_PATTERNS } from '@/shared/config/routes';

import { AppLayout } from './AppLayout';
import { RequireAuth, RequireGuest } from './guards';

const ItemsListPage = lazy(() =>
  import('@/pages/items-list').then((module) => ({ default: module.ItemsListPage })),
);
const ItemDetailPage = lazy(() =>
  import('@/pages/item-detail').then((module) => ({ default: module.ItemDetailPage })),
);
const ItemCreatePage = lazy(() =>
  import('@/pages/item-create').then((module) => ({ default: module.ItemCreatePage })),
);
const ItemEditPage = lazy(() =>
  import('@/pages/item-edit').then((module) => ({ default: module.ItemEditPage })),
);
const MyItemsPage = lazy(() =>
  import('@/pages/my-items').then((module) => ({ default: module.MyItemsPage })),
);
const PetPage = lazy(() => import('@/pages/pet').then((module) => ({ default: module.PetPage })));
const RewardsPage = lazy(() =>
  import('@/pages/rewards').then((module) => ({ default: module.RewardsPage })),
);
const LeaderboardPage = lazy(() =>
  import('@/pages/leaderboard').then((module) => ({ default: module.LeaderboardPage })),
);
const LoginPage = lazy(() =>
  import('@/pages/login').then((module) => ({ default: module.LoginPage })),
);
const RegisterPage = lazy(() =>
  import('@/pages/register').then((module) => ({ default: module.RegisterPage })),
);
const ProfilePage = lazy(() =>
  import('@/pages/profile').then((module) => ({ default: module.ProfilePage })),
);
const NotFoundPage = lazy(() =>
  import('@/pages/not-found').then((module) => ({ default: module.NotFoundPage })),
);

export const router = createBrowserRouter([
  {
    path: ROUTES.home,
    element: <AppLayout />,
    children: [
      { index: true, element: <Navigate to={ROUTES.items} replace /> },
      { path: ROUTES.items, element: <ItemsListPage /> },
      { path: ROUTE_PATTERNS.itemDetail, element: <ItemDetailPage /> },
      {
        path: ROUTES.itemCreate,
        element: (
          <RequireAuth>
            <ItemCreatePage />
          </RequireAuth>
        ),
      },
      {
        path: ROUTE_PATTERNS.itemEdit,
        element: (
          <RequireAuth>
            <ItemEditPage />
          </RequireAuth>
        ),
      },
      {
        path: ROUTES.myItems,
        element: (
          <RequireAuth>
            <MyItemsPage />
          </RequireAuth>
        ),
      },
      {
        path: ROUTES.pet,
        element: (
          <RequireAuth>
            <PetPage />
          </RequireAuth>
        ),
      },
      {
        path: ROUTES.rewards,
        element: (
          <RequireAuth>
            <RewardsPage />
          </RequireAuth>
        ),
      },
      {
        path: ROUTES.leaderboard,
        element: (
          <RequireAuth>
            <LeaderboardPage />
          </RequireAuth>
        ),
      },
      {
        path: ROUTES.profile,
        element: (
          <RequireAuth>
            <ProfilePage />
          </RequireAuth>
        ),
      },
      {
        path: ROUTES.login,
        element: (
          <RequireGuest>
            <LoginPage />
          </RequireGuest>
        ),
      },
      {
        path: ROUTES.register,
        element: (
          <RequireGuest>
            <RegisterPage />
          </RequireGuest>
        ),
      },
      { path: '*', element: <NotFoundPage /> },
    ],
  },
]);
