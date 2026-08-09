import { describe, expect, it } from 'vitest'
import { ApiError } from '#/api/api-error'
import { toServerFieldName, toServerMessage } from './form-errors'

const t = (key: string, options?: Record<string, unknown>) =>
  options ? `${key}:${JSON.stringify(options)}` : key

describe('toServerMessage', () => {
  it('translates a short password using the boundary from the schema', () => {
    const message = toServerMessage(
      new ApiError({
        kind: 'validation_error',
        message: 'value is shorter than minimum: 8',
        status: 400,
        field: 'password',
      }),
      t,
    )

    expect(message).toBe('validation:passwordTooShort:{"count":8}')
  })

  it('translates a long password', () => {
    const message = toServerMessage(
      new ApiError({
        kind: 'validation_error',
        message: 'value is longer than maximum: 72',
        status: 400,
        field: 'password',
      }),
      t,
    )

    expect(message).toBe('validation:passwordTooLong:{"count":72}')
  })

  it('never leaks the raw english server text to the user', () => {
    const message = toServerMessage(
      new ApiError({
        kind: 'validation_error',
        message: 'value is shorter than minimum: 8',
        status: 400,
        field: 'password',
      }),
      t,
    )

    expect(message).not.toContain('minimum')
  })

  it('maps a 409 to the taken email message', () => {
    expect(
      toServerMessage(
        new ApiError({
          kind: 'conflict',
          message: 'state conflict',
          status: 409,
        }),
        t,
      ),
    ).toBe('auth:errors.emailTaken')
  })

  it('maps a 401 to invalid credentials', () => {
    expect(
      toServerMessage(
        new ApiError({
          kind: 'unauthorized',
          message: 'authentication required',
          status: 401,
        }),
        t,
      ),
    ).toBe('auth:errors.invalidCredentials')
  })

  it('maps transport failures to a readable message', () => {
    expect(
      toServerMessage(new ApiError({ kind: 'network', message: 'x' }), t),
    ).toBe('errors:network')
  })

  it('returns nothing without an error', () => {
    expect(toServerMessage(null, t)).toBeNull()
  })
})

describe('toServerFieldName', () => {
  it('points a password validation error at the password field', () => {
    expect(
      toServerFieldName(
        new ApiError({
          kind: 'validation_error',
          message: 'too short',
          status: 400,
          field: 'password',
        }),
      ),
    ).toBe('password')
  })

  it('maps the snake case full_name to the form field name', () => {
    expect(
      toServerFieldName(
        new ApiError({
          kind: 'validation_error',
          message: 'field is required',
          status: 400,
          field: 'full_name',
        }),
      ),
    ).toBe('fullName')
  })

  it('points a duplicate email conflict at the email field', () => {
    expect(
      toServerFieldName(
        new ApiError({
          kind: 'conflict',
          message: 'state conflict',
          status: 409,
        }),
      ),
    ).toBe('email')
  })

  it('does not blame a field for a 401', () => {
    expect(
      toServerFieldName(
        new ApiError({ kind: 'unauthorized', message: 'x', status: 401 }),
      ),
    ).toBeNull()
  })
})
