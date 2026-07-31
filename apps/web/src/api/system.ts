import {
  selectDirectory,
  type ErrorResponse,
  type SelectDirectoryResponse,
} from '@corvus-studio/api-client';

export class DirectoryPickerApiError extends Error {
  readonly code: string;
  readonly recoverable: boolean;

  constructor(error: ErrorResponse['error']) {
    super(error.message);
    this.name = 'DirectoryPickerApiError';
    this.code = error.code;
    this.recoverable = error.recoverable;
  }
}

const requestBaseURL = window.location.origin;

export const systemApi = {
  async selectProjectDirectory(): Promise<string | null> {
    try {
      const { data } = await selectDirectory({
        baseUrl: requestBaseURL,
        body: { purpose: 'project_location' },
        throwOnError: true,
      });
      return selectedPath(data);
    } catch (error) {
      throw normalizeDirectoryPickerError(error);
    }
  },
};

function selectedPath(response: SelectDirectoryResponse): string | null {
  if (!response.selected) {
    return null;
  }
  if (response.path === null || response.path.length === 0) {
    throw new Error('Corvus Core returned an empty directory selection');
  }
  return response.path;
}

function normalizeDirectoryPickerError(error: unknown): Error {
  if (isErrorResponse(error)) {
    return new DirectoryPickerApiError(error.error);
  }
  if (error instanceof Error) {
    return error;
  }
  return new Error('Corvus Core could not open the directory picker');
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
