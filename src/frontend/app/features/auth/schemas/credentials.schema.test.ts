import { describe, expect, it } from 'vitest'
import { SignInSchema } from './sign-in.schema'
import { SignUpSchema } from './sign-up.schema'

const validEmail = 'demo@example.com'

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

describe('credential validation parity', () => {
  it('accepts the same short password on sign in and sign up', () => {
    const password = 'a'

    expect(signInErrors({ email: validEmail, password })).toEqual([])
    expect(
      signUpErrors({ email: validEmail, password, fullName: 'Иван' }),
    ).toEqual([])
  })

  it('imposes no minimum password length on either form', () => {
    for (const password of ['a', 'ab', '1234', 'short']) {
      expect(signInErrors({ email: validEmail, password })).toEqual([])
      expect(
        signUpErrors({ email: validEmail, password, fullName: 'Иван' }),
      ).toEqual([])
    }
  })

  it('rejects an empty password on both forms', () => {
    expect(signInErrors({ email: validEmail, password: '' })).toEqual([
      'password',
    ])
    expect(
      signUpErrors({ email: validEmail, password: '', fullName: 'Иван' }),
    ).toEqual(['password'])
  })

  it('applies the same email rule to both forms', () => {
    for (const email of ['not-an-email', 'a@b', 'a b@c.dev', '']) {
      expect(signInErrors({ email, password: 'pw' })).toContain('email')
      expect(
        signUpErrors({ email, password: 'pw', fullName: 'Иван' }),
      ).toContain('email')
    }
  })

  it('requires a non-empty full name only on sign up', () => {
    expect(
      signUpErrors({ email: validEmail, password: 'pw', fullName: '   ' }),
    ).toEqual(['fullName'])
    expect(signInErrors({ email: validEmail, password: 'pw' })).toEqual([])
  })

  it('uses translation keys for messages', () => {
    const result = SignInSchema.safeParse({ email: 'nope', password: '' })

    expect(result.success).toBe(false)
    if (result.success) return

    const messages = result.error.issues.map((issue) => issue.message)
    expect(messages).toContain('validation:email')
    expect(messages).toContain('validation:required')
  })
})
