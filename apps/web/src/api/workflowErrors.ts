import type { ErrorResponse } from '@corvus-studio/api-client';

export class WorkflowApiError extends Error {
  readonly code: string;
  readonly recoverable: boolean;

  constructor(error: ErrorResponse['error']) {
    super(error.message);
    this.name = 'WorkflowApiError';
    this.code = error.code;
    this.recoverable = error.recoverable;
  }
}

export function normalizeWorkflowError(error: unknown): Error {
  if (isErrorResponse(error)) {
    return new WorkflowApiError(error.error);
  }
  if (error instanceof Error) {
    return error;
  }
  return new Error('Corvus Core could not be reached');
}

function isErrorResponse(value: unknown): value is ErrorResponse {
  if (typeof value !== 'object' || value === null || !('error' in value)) {
    return false;
  }
  const error = value.error;
  return (
    typeof error === 'object' &&
    error !== null &&
    'code' in error &&
    typeof error.code === 'string' &&
    'message' in error &&
    typeof error.message === 'string' &&
    'recoverable' in error &&
    typeof error.recoverable === 'boolean'
  );
}
