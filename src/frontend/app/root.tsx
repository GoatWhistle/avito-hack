import {
  Links,
  type LinksFunction,
  Meta,
  Outlet,
  Scripts,
  ScrollRestoration,
} from 'react-router'
import { Providers } from '#/providers'
import { AppLayout } from '#/features/layout'
import { Devtools } from '#/devtools'
import type { PropsWithChildren } from 'react'

import './app.css'

export const links: LinksFunction = () => []

export function Layout({ children }: PropsWithChildren) {
  return (
    <html lang="ru">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <Meta />
        <Links />
      </head>
      <body>
        <Providers>
          {children}
          <Devtools />
        </Providers>
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  )
}

export default function App() {
  return (
    <AppLayout>
      <Outlet />
    </AppLayout>
  )
}
