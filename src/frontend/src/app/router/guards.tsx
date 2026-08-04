import { useUnit } from 'effector-react';
import type { PropsWithChildren } from 'react';
import { Navigate, useLocation } from 'react-router-dom';

import { $isAuthenticated, $isSessionPending } from '@/entities/session';
import { ROUTES } from '@/shared/config/routes';
import { PageSkeleton } from '@/shared/ui';

export function RequireAuth({ children }: PropsWithChildren) {
  const [isAuthenticated, isPending] = useUnit([$isAuthenticated, $isSessionPending]);
  const location = useLocation();

  if (isPending) {
    return <PageSkeleton />;
  }

  if (!isAuthenticated) {
    return <Navigate to={ROUTES.login} state={{ from: location.pathname }} replace />;
  }

  return children;
}

export function RequireGuest({ children }: PropsWithChildren) {
  const [isAuthenticated, isPending] = useUnit([$isAuthenticated, $isSessionPending]);

  if (isPending) {
    return <PageSkeleton />;
  }

  if (isAuthenticated) {
    return <Navigate to={ROUTES.myItems} replace />;
  }

  return children;
}
