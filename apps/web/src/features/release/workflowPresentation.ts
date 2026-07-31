import type {
  ChecklistCategory,
  ChecklistRequirementLevel,
  ChecklistSource,
  ChecklistStatus,
} from '../../api/checklist';
import type { ReleaseStatus } from '../../api/releases';

export const releaseStatusLabels: Record<ReleaseStatus, string> = {
  draft: 'Draft',
  preparing: 'Preparing',
  needs_attention: 'Needs attention',
  ready_for_review: 'Ready for review',
  ready: 'Ready',
  submitted: 'Submitted',
};

export const checklistStatusLabels: Record<ChecklistStatus, string> = {
  not_started: 'Not started',
  in_progress: 'In progress',
  needs_review: 'Needs review',
  done: 'Done',
  blocked: 'Blocked',
  not_applicable: 'Not applicable',
};

export const sourceLabels: Record<ChecklistSource, string> = {
  platform_template: 'Steam',
  corvus_template: 'Corvus',
  user: 'Custom',
  agent: 'Agent',
};

export const requirementLevelLabels: Record<ChecklistRequirementLevel, string> = {
  required: 'Required',
  recommended: 'Recommended',
};

export const categoryLabels: Record<ChecklistCategory, string> = {
  setup: 'Setup',
  store_copy: 'Store copy',
  branding: 'Branding',
  media: 'Media',
  compliance: 'Compliance',
  timeline: 'Timeline',
  positioning: 'Positioning',
  localization: 'Localization',
  review: 'Review',
};

export const releaseTransitions: Record<ReleaseStatus, ReleaseStatus[]> = {
  draft: ['preparing'],
  preparing: ['needs_attention', 'ready_for_review'],
  needs_attention: ['preparing'],
  ready_for_review: ['preparing', 'needs_attention', 'ready'],
  ready: ['preparing', 'needs_attention', 'submitted'],
  submitted: ['needs_attention'],
};
