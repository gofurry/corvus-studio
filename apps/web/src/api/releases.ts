import {
  createRelease as createReleaseRequest,
  getRelease as getReleaseRequest,
  getReleaseTemplate as getReleaseTemplateRequest,
  listProjectReleases as listProjectReleasesRequest,
  transitionRelease as transitionReleaseRequest,
  type CreateReleaseRequest,
  type ReleaseGoal,
  type ReleaseStatus,
  type ReleaseTemplate,
} from '@corvus-studio/api-client';

import { normalizeWorkflowError } from './workflowErrors';

export type { CreateReleaseRequest, ReleaseGoal, ReleaseStatus, ReleaseTemplate };

const requestBaseURL = window.location.origin;

export const releasesApi = {
  async getTemplate(templateKey: string): Promise<ReleaseTemplate> {
    try {
      const { data } = await getReleaseTemplateRequest({
        baseUrl: requestBaseURL,
        path: { template_key: templateKey },
        throwOnError: true,
      });
      return data;
    } catch (error) {
      throw normalizeWorkflowError(error);
    }
  },

  async listForProject(projectId: string): Promise<ReleaseGoal[]> {
    try {
      const { data } = await listProjectReleasesRequest({
        baseUrl: requestBaseURL,
        path: { project_id: projectId },
        throwOnError: true,
      });
      return data.items;
    } catch (error) {
      throw normalizeWorkflowError(error);
    }
  },

  async create(input: CreateReleaseRequest): Promise<ReleaseGoal> {
    try {
      const { data } = await createReleaseRequest({
        baseUrl: requestBaseURL,
        body: input,
        throwOnError: true,
      });
      return data;
    } catch (error) {
      throw normalizeWorkflowError(error);
    }
  },

  async get(releaseId: string): Promise<ReleaseGoal> {
    try {
      const { data } = await getReleaseRequest({
        baseUrl: requestBaseURL,
        path: { release_id: releaseId },
        throwOnError: true,
      });
      return data;
    } catch (error) {
      throw normalizeWorkflowError(error);
    }
  },

  async transition(releaseId: string, status: ReleaseStatus): Promise<ReleaseGoal> {
    try {
      const { data } = await transitionReleaseRequest({
        baseUrl: requestBaseURL,
        path: { release_id: releaseId },
        body: { status },
        throwOnError: true,
      });
      return data;
    } catch (error) {
      throw normalizeWorkflowError(error);
    }
  },
};
