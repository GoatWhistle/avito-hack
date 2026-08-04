export const ERROR_CODES = [
  'not_found',
  'forbidden',
  'unauthorized',
  'conflict',
  'validation_error',
  'bad_request',
  'internal_error',
  'network_error',
] as const;

export type ErrorCode = (typeof ERROR_CODES)[number];

export interface ApiErrorBody {
  error: {
    code: string;
    message: string;
    field?: string;
    request_id?: string;
  };
}

export interface ApiListResponse<T> {
  items: T[];
  next_cursor?: string;
}

export class ApiError extends Error {
  readonly code: ErrorCode;
  readonly field: string | undefined;
  readonly requestId: string | undefined;
  readonly status: number;

  constructor(params: {
    code: ErrorCode;
    message: string;
    status: number;
    field?: string;
    requestId?: string;
  }) {
    super(params.message);
    this.name = 'ApiError';
    this.code = params.code;
    this.status = params.status;
    this.field = params.field;
    this.requestId = params.requestId;
  }
}

export function isErrorCode(value: string): value is ErrorCode {
  return (ERROR_CODES as readonly string[]).includes(value);
}
