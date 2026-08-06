import { TanStackDevtools } from '@tanstack/react-devtools'
import { ReactQueryDevtoolsPanel } from '@tanstack/react-query-devtools'
import { formDevtoolsPlugin } from '@tanstack/react-form-devtools'

export function Devtools() {
  return (
    <TanStackDevtools
      plugins={[
        formDevtoolsPlugin(),
        {
          name: 'TanStack Query',
          render: <ReactQueryDevtoolsPanel />,
        },
      ]}
    />
  )
}
