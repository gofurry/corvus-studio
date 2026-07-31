import type { ChecklistFilters } from '../../api/checklist';

export const releaseQueryKeys = {
  template: (templateKey: string) => ['release-template', templateKey] as const,
  project: (projectId: string) => ['releases', 'project', projectId] as const,
  detail: (releaseId: string) => ['releases', releaseId] as const,
};

export const checklistQueryKeys = {
  list: (releaseId: string, filters: ChecklistFilters) =>
    ['checklist', releaseId, filters] as const,
  detail: (itemId: string) => ['checklist-item', itemId] as const,
};
