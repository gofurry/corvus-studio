import {
  createProject as createProjectRequest,
  getProject as getProjectRequest,
  listProjects as listProjectsRequest,
  type CreateProjectRequest,
  type ErrorResponse,
  type Project,
  type ProjectStage,
  type ProjectStatus,
} from '@corvus-studio/api-client';

export type { CreateProjectRequest, Project, ProjectStage, ProjectStatus };

export class ProjectApiError extends Error {
  readonly code: string;
  readonly recoverable: boolean;

  constructor(error: ErrorResponse['error']) {
    super(error.message);
    this.name = 'ProjectApiError';
    this.code = error.code;
    this.recoverable = error.recoverable;
  }
}

const requestBaseURL = window.location.origin;

export const projectsApi = {
  async list(): Promise<Project[]> {
    try {
      const { data } = await listProjectsRequest({
        baseUrl: requestBaseURL,
        throwOnError: true,
      });
      return data.items;
    } catch (error) {
      throw normalizeProjectError(error);
    }
  },

  async create(input: CreateProjectRequest): Promise<Project> {
    try {
      const { data } = await createProjectRequest({
        baseUrl: requestBaseURL,
        body: input,
        throwOnError: true,
      });
      return data;
    } catch (error) {
      throw normalizeProjectError(error);
    }
  },

  async get(projectId: string): Promise<Project> {
    try {
      const { data } = await getProjectRequest({
        baseUrl: requestBaseURL,
        path: { project_id: projectId },
        throwOnError: true,
      });
      return data;
    } catch (error) {
      throw normalizeProjectError(error);
    }
  },
};

function normalizeProjectError(error: unknown): Error {
  if (isErrorResponse(error)) {
    return new ProjectApiError(error.error);
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
