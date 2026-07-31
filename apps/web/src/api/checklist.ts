import {
  createChecklistItem as createChecklistItemRequest,
  getChecklistItem as getChecklistItemRequest,
  listChecklistItems as listChecklistItemsRequest,
  transitionChecklistItem as transitionChecklistItemRequest,
  type ChecklistCategory,
  type ChecklistItem,
  type ChecklistRequirementLevel,
  type ChecklistSource,
  type ChecklistStatus,
  type CreateChecklistItemRequest,
} from '@corvus-studio/api-client';

import { normalizeWorkflowError } from './workflowErrors';

export type {
  ChecklistCategory,
  ChecklistItem,
  ChecklistRequirementLevel,
  ChecklistSource,
  ChecklistStatus,
  CreateChecklistItemRequest,
};

export type ChecklistFilters = {
  status?: ChecklistStatus;
  source?: ChecklistSource;
  category?: ChecklistCategory;
};

const requestBaseURL = window.location.origin;

export const checklistApi = {
  async list(releaseId: string, filters: ChecklistFilters = {}): Promise<ChecklistItem[]> {
    try {
      const { data } = await listChecklistItemsRequest({
        baseUrl: requestBaseURL,
        path: { release_id: releaseId },
        query: filters,
        throwOnError: true,
      });
      return data.items;
    } catch (error) {
      throw normalizeWorkflowError(error);
    }
  },

  async create(input: CreateChecklistItemRequest): Promise<ChecklistItem> {
    try {
      const { data } = await createChecklistItemRequest({
        baseUrl: requestBaseURL,
        body: input,
        throwOnError: true,
      });
      return data;
    } catch (error) {
      throw normalizeWorkflowError(error);
    }
  },

  async get(itemId: string): Promise<ChecklistItem> {
    try {
      const { data } = await getChecklistItemRequest({
        baseUrl: requestBaseURL,
        path: { item_id: itemId },
        throwOnError: true,
      });
      return data;
    } catch (error) {
      throw normalizeWorkflowError(error);
    }
  },

  async transition(itemId: string, status: ChecklistStatus): Promise<ChecklistItem> {
    try {
      const { data } = await transitionChecklistItemRequest({
        baseUrl: requestBaseURL,
        path: { item_id: itemId },
        body: { status },
        throwOnError: true,
      });
      return data;
    } catch (error) {
      throw normalizeWorkflowError(error);
    }
  },
};
