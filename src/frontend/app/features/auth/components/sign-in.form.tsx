import { useSignInForm } from '#/features/auth/forms'
import {
  Button,
  CardContent,
  CardFooter,
  Field,
  FieldGroup,
  FieldLabel,
  Input,
} from '#/components/ui'

export function SignInForm() {
  const form = useSignInForm()

  return (
    <form
      onSubmit={async (e) => {
        e.preventDefault()
        e.stopPropagation()
        await form.handleSubmit()
      }}
    >
      <CardContent>
        <FieldGroup>
          <form.Field
            name="email"
            children={(field) => {
              return (
                <Field>
                  <FieldLabel htmlFor={field.name}>
                    Электронная почта
                  </FieldLabel>
                  <Input
                    id={field.name}
                    name={field.name}
                    value={field.state.value}
                    onChange={(e) => field.handleChange(e.target.value)}
                  />
                </Field>
              )
            }}
          />

          <form.Field
            name="password"
            children={(field) => {
              return (
                <Field>
                  <FieldLabel htmlFor={field.name}>Пароль</FieldLabel>
                  <Input
                    id={field.name}
                    name={field.name}
                    value={field.state.value}
                    onChange={(e) => field.handleChange(e.target.value)}
                  />
                </Field>
              )
            }}
          />
        </FieldGroup>
      </CardContent>

      <CardFooter>
        <form.Subscribe
          selector={(state) => [state.canSubmit]}
          children={([canSubmit]) => (
            <Button type="submit" disabled={!canSubmit}>
              Войти
            </Button>
          )}
        />
      </CardFooter>
    </form>
  )
}
