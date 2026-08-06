import { useSignUpForm } from '#/features/auth/forms'
import {
  Button,
  CardContent,
  CardFooter,
  Field,
  FieldGroup,
  FieldLabel,
  Input,
} from '#/components/ui'

export function SignUpForm() {
  const form = useSignUpForm()

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

          <form.Field
            name="fullName"
            children={(field) => {
              return (
                <Field>
                  <FieldLabel htmlFor={field.name}>Полное имя</FieldLabel>
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
              Зарегистрироваться
            </Button>
          )}
        />
      </CardFooter>
    </form>
  )
}
