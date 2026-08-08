import { Providers } from '#/providers/Providers'
import { Toaster } from '#/shared/components/ui/toast'
import { Devtools } from '#/shared/devtools/Devtools'
import { usePetWebSocket } from '#/shared/hooks/usePetWebSocket'
import './app.css'
import {
  Links,
  type LinksFunction,
  Meta,
  Outlet,
  Scripts,
  ScrollRestoration,
} from 'react-router'

export const links: LinksFunction = () => [
  { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
  {
    rel: 'preconnect',
    href: 'https://fonts.gstatic.com',
    crossOrigin: 'anonymous',
  },
  {
    rel: 'stylesheet',
    href: 'https://fonts.googleapis.com/css2?family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&display=swap',
  },
]

export function Layout() {
  return (
    <html lang="ru">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <Meta />
        <Links />
      </head>
      <body className="min-h-svh bg-background font-sans text-foreground antialiased selection:bg-primary/20 selection:text-primary relative overflow-x-hidden">
        <div
          className="absolute inset-0 bg-[radial-gradient(ellipse_100%_60%_at_50%_-10%,rgba(255,200,140,0.18),transparent_70%)] pointer-events-none"
          aria-hidden
        />

        <div
          className="absolute top-20 -left-40 size-[600px] rounded-full bg-primary/15 blur-[120px] pointer-events-none animate-blob"
          aria-hidden
        />

        <div
          className="absolute top-[40%] right-0 size-[500px] rounded-full bg-[#97CF26]/15 blur-[140px] pointer-events-none animate-blob animation-delay-2000 dark:bg-[#97CF26]/10"
          aria-hidden
        />

        <div
          className="absolute bottom-20 left-[30%] size-[550px] rounded-full bg-[#A169F7]/15 blur-[130px] pointer-events-none animate-blob animation-delay-4000 dark:bg-[#A169F7]/10"
          aria-hidden
        />

        <div
          className="absolute top-1/2 right-1/4 size-[400px] rounded-full bg-[#FF6163]/10 blur-[120px] pointer-events-none animate-blob animation-delay-3000 dark:bg-[#FF6163]/8"
          aria-hidden
        />

        <div
          className="absolute inset-0 bg-[radial-gradient(circle_at_1px_1px,rgba(0,0,0,0.08)_1px,transparent_0)] bg-[size:24px_24px] pointer-events-none dark:bg-[radial-gradient(circle_at_1px_1px,rgba(255,255,255,0.05)_1px,transparent_0)]"
          aria-hidden
        />

        <div
          className="absolute top-0 left-0 size-[800px] bg-[linear-gradient(135deg,rgba(255,255,255,0.4),transparent_50%)] pointer-events-none dark:bg-[linear-gradient(135deg,rgba(255,255,255,0.05),transparent_50%)]"
          aria-hidden
        />

        <div className="relative z-10 flex min-h-svh flex-col">
          <Providers>
            <Outlet />
            <Toaster />
            <Devtools />
          </Providers>
        </div>

        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  )
}

export default function App() {
  usePetWebSocket()

  return <Outlet />
}
