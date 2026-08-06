import { TanStackDevtools } from '@tanstack/react-devtools'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { FormDevtoolsPanel } from '@tanstack/react-form-devtools'

export function Devtools() {
  return (
    <TanStackDevtools
      plugins={[
        {
          name: 'TanStack Query',
          render: <ReactQueryDevtools />,
        },
        {
          name: 'TanStack Form',
          render: <FormDevtoolsPanel />,
        },
      ]}
    />
  )
}
