import {
  isRouteErrorResponse,
  Links,
  type LinksFunction,
  Meta,
  type MetaFunction,
  Outlet,
  Scripts,
  ScrollRestoration,
  useLocation,
  useRouteError,
} from 'react-router'
import { Providers } from '#/providers'
import { AppLayout, Aurora } from '#/features/layout'
import { CrashScreen, NotFoundScreen } from '#/features/errors'
import { fallbackLocale, resources } from '#/i18n/resources'
import { themeColors } from '#/tokens'
import { routeTransitionKey } from '#/route-transition-key'
import type { PropsWithChildren } from 'react'

import './app.css'

export const meta: MetaFunction = () => {
  const { app } = resources[fallbackLocale].common

  return [
    { title: app.metaTitle },
    { name: 'description', content: app.metaDescription },
    { property: 'og:type', content: 'website' },
    { property: 'og:site_name', content: app.name },
    { property: 'og:title', content: app.metaTitle },
    { property: 'og:description', content: app.metaDescription },
    { property: 'og:image', content: '/og-image.png' },
    { property: 'og:image:width', content: '1200' },
    { property: 'og:image:height', content: '630' },
    { property: 'og:locale', content: 'ru_RU' },
    { name: 'twitter:card', content: 'summary_large_image' },
    { name: 'twitter:title', content: app.metaTitle },
    { name: 'twitter:description', content: app.metaDescription },
    { name: 'twitter:image', content: '/og-image.png' },
  ]
}

export const links: LinksFunction = () => [
  { rel: 'icon', href: '/favicon.ico', sizes: '32x32' },
  { rel: 'icon', href: '/favicon.svg', type: 'image/svg+xml' },
  { rel: 'apple-touch-icon', href: '/apple-touch-icon.png', sizes: '180x180' },
  { rel: 'manifest', href: '/site.webmanifest' },
]

export function Layout({ children }: PropsWithChildren) {
  return (
    <html lang="ru">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta
          name="theme-color"
          content={themeColors.light}
          media="(prefers-color-scheme: light)"
        />
        <meta
          name="theme-color"
          content={themeColors.dark}
          media="(prefers-color-scheme: dark)"
        />
        <Meta />
        <Links />
      </head>
      <body>
        <Aurora />
        <Providers>{children}</Providers>
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  )
}

function RouteTransition() {
  const { pathname } = useLocation()

  return (
    <div key={routeTransitionKey(pathname)} className="route-transition">
      <Outlet />
    </div>
  )
}

export default function App() {
  return (
    <AppLayout>
      <RouteTransition />
    </AppLayout>
  )
}

export function ErrorBoundary() {
  const error = useRouteError()
  const isNotFound = isRouteErrorResponse(error) && error.status === 404

  return (
    <AppLayout>{isNotFound ? <NotFoundScreen /> : <CrashScreen />}</AppLayout>
  )
}
