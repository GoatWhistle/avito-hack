import { describe, expect, it } from 'vitest'
import { PASSWORD_MAX_LENGTH, PASSWORD_MIN_LENGTH } from './credentials.schema'
import { SignInSchema } from './sign-in.schema'
import { SignUpSchema } from './sign-up.schema'

const validEmail = 'demo@example.com'
const validPassword = 'demo1234'

const signInErrors = (value: { email: string; password: string }) => {
  const result = SignInSchema.safeParse(value)

  return result.success ? [] : result.error.issues.map((issue) => issue.path[0])
}

const signUpErrors = (value: {
  email: string
  password: string
  fullName: string
}) => {
  const result = SignUpSchema.safeParse(value)

  return result.success ? [] : result.error.issues.map((issue) => issue.path[0])
}

describe('credential validation', () => {
  it('accepts a password at both length boundaries on sign up', () => {
    for (const length of [PASSWORD_MIN_LENGTH, PASSWORD_MAX_LENGTH]) {
      const password = 'a'.repeat(length)

      expect(
        signUpErrors({ email: validEmail, password, fullName: 'Иван' }),
      ).toEqual([])
    }
  })

  it('rejects a password shorter than the backend minimum on sign up', () => {
    for (const length of [1, 4, PASSWORD_MIN_LENGTH - 1]) {
      expect(
        signUpErrors({
          email: validEmail,
          password: 'a'.repeat(length),
          fullName: 'Иван',
        }),
      ).toContain('password')
    }
  })

  it('rejects a password longer than the backend maximum on sign up', () => {
    expect(
      signUpErrors({
        email: validEmail,
        password: 'a'.repeat(PASSWORD_MAX_LENGTH + 1),
        fullName: 'Иван',
      }),
    ).toContain('password')
  })

  it('keeps sign in permissive so existing accounts can still authenticate', () => {
    for (const password of ['a', 'short', validPassword]) {
      expect(signInErrors({ email: validEmail, password })).toEqual([])
    }
  })

  it('rejects an empty password on both forms', () => {
    expect(signInErrors({ email: validEmail, password: '' })).toEqual([
      'password',
    ])
    expect(
      signUpErrors({ email: validEmail, password: '', fullName: 'Иван' }),
    ).toContain('password')
  })

  it('applies the same email rule to both forms', () => {
    for (const email of ['not-an-email', 'a@b', 'a b@c.dev', '']) {
      expect(signInErrors({ email, password: validPassword })).toContain(
        'email',
      )
      expect(
        signUpErrors({ email, password: validPassword, fullName: 'Иван' }),
      ).toContain('email')
    }
  })

  it('requires a non-empty full name only on sign up', () => {
    expect(
      signUpErrors({
        email: validEmail,
        password: validPassword,
        fullName: '   ',
      }),
    ).toEqual(['fullName'])
    expect(
      signInErrors({ email: validEmail, password: validPassword }),
    ).toEqual([])
  })

  it('uses translation keys for messages', () => {
    const result = SignUpSchema.safeParse({
      email: 'nope',
      password: 'short',
      fullName: '',
    })

    expect(result.success).toBe(false)
    if (result.success) return

    const messages = result.error.issues.map((issue) => issue.message)
    expect(messages).toContain('validation:email')
    expect(messages).toContain('validation:passwordTooShort')
    expect(messages).toContain('validation:required')
  })
})
