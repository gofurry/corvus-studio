import type { ProjectStage, ProjectStatus } from '../../api/projects';

export const projectStageLabels: Record<ProjectStage, string> = {
  concept: 'Concept',
  development: 'Development',
  release_preparation: 'Release preparation',
  released: 'Released',
};

export const projectStatusLabels: Record<ProjectStatus, string> = {
  active: 'Active',
  archived: 'Archived',
};

export function formatProjectDate(value: string): string {
  return new Intl.DateTimeFormat(undefined, {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(value));
}
