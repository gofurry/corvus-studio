import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import App from '../../App';
import { checklistApi, type ChecklistItem } from '../../api/checklist';
import { projectsApi, type Project } from '../../api/projects';
import { releasesApi, type ReleaseGoal, type ReleaseTemplate } from '../../api/releases';

vi.mock('../../api/projects', () => ({
  projectsApi: {
    list: vi.fn(),
    create: vi.fn(),
    get: vi.fn(),
  },
}));

vi.mock('../../api/releases', () => ({
  releasesApi: {
    getTemplate: vi.fn(),
    listForProject: vi.fn(),
    create: vi.fn(),
    get: vi.fn(),
    transition: vi.fn(),
  },
}));

vi.mock('../../api/checklist', () => ({
  checklistApi: {
    list: vi.fn(),
    create: vi.fn(),
    get: vi.fn(),
    transition: vi.fn(),
  },
}));

const mockedProjectsApi = vi.mocked(projectsApi);
const mockedReleasesApi = vi.mocked(releasesApi);
const mockedChecklistApi = vi.mocked(checklistApi);

describe('Release and Checklist workspace', () => {
  beforeEach(() => {
    vi.resetAllMocks();
    mockedProjectsApi.get.mockResolvedValue(testProject());
    mockedReleasesApi.getTemplate.mockResolvedValue(testTemplate());
    mockedChecklistApi.list.mockResolvedValue([]);
  });

  afterEach(() => {
    cleanup();
  });

  it('previews and confirms one-time Release workspace generation', async () => {
    const release = testRelease();
    mockedReleasesApi.listForProject.mockResolvedValue([]);
    mockedReleasesApi.create.mockResolvedValue(release);
    renderApp(`/projects/${testProjectId}/release`);

    expect(await screen.findByText('Template version')).toBeInTheDocument();
    expect(screen.getByText('Steam tasks')).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Generate release workspace' }));
    fireEvent.click(await screen.findByRole('button', { name: 'Generate workspace' }));

    await waitFor(() => {
      expect(mockedReleasesApi.create).toHaveBeenCalledWith({
        project_id: testProjectId,
        goal_type: 'steam_coming_soon',
        template_key: 'steam-coming-soon',
        template_version: '1.0.0',
      });
    });
    expect(await screen.findByRole('heading', { name: 'Release checklist' })).toBeInTheDocument();
  });

  it('restores an existing Release Goal without reloading the creation template', async () => {
    const release = testRelease();
    mockedReleasesApi.listForProject.mockResolvedValue([release]);
    renderApp(`/projects/${testProjectId}/release`);

    expect(await screen.findByText('1/12')).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Open checklist' })).toBeInTheDocument();
    expect(mockedReleasesApi.getTemplate).not.toHaveBeenCalled();
  });

  it('restores filters and selected task from a direct Checklist URL', async () => {
    const release = testRelease();
    const item = testItem();
    mockedReleasesApi.listForProject.mockResolvedValue([release]);
    mockedChecklistApi.list.mockResolvedValue([item]);
    mockedChecklistApi.get.mockResolvedValue(item);
    renderApp(
      `/projects/${testProjectId}/checklist?status=in_progress&source=user&category=review&item=${item.id}`,
    );

    expect(await screen.findByRole('heading', { name: item.title })).toBeInTheDocument();
    expect(mockedChecklistApi.list).toHaveBeenCalledWith(release.id, {
      status: 'in_progress',
      source: 'user',
      category: 'review',
    });
    expect(mockedChecklistApi.get).toHaveBeenCalledWith(item.id);
  });

  it('adds a user-owned custom task', async () => {
    const release = testRelease();
    const created = testItem();
    mockedReleasesApi.listForProject.mockResolvedValue([release]);
    mockedChecklistApi.create.mockResolvedValue(created);
    mockedChecklistApi.get.mockResolvedValue(created);
    renderApp(`/projects/${testProjectId}/checklist`);

    fireEvent.click(await screen.findByRole('button', { name: 'Add custom task' }));
    fireEvent.change(screen.getByLabelText('Task title'), {
      target: { value: created.title },
    });
    fireEvent.click(screen.getByRole('button', { name: 'Add task' }));

    await waitFor(() => {
      expect(mockedChecklistApi.create).toHaveBeenCalledWith(
        expect.objectContaining({
          release_goal_id: release.id,
          title: created.title,
          category: 'review',
          requirement_level: 'recommended',
        }),
      );
    });
  });

  it('changes Checklist task status from the details drawer', async () => {
    const release = testRelease();
    const item = testItem();
    mockedReleasesApi.listForProject.mockResolvedValue([release]);
    mockedChecklistApi.list.mockResolvedValue([item]);
    mockedChecklistApi.get.mockResolvedValue(item);
    mockedChecklistApi.transition.mockResolvedValue({ ...item, status: 'done' });
    renderApp(`/projects/${testProjectId}/checklist?item=${item.id}`);

    fireEvent.click(await screen.findByRole('button', { name: 'Set task status to Done' }));

    await waitFor(() => {
      expect(mockedChecklistApi.transition).toHaveBeenCalledWith(item.id, 'done');
    });
  });
});

function renderApp(initialRoute: string) {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[initialRoute]}>
        <App />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

const testProjectId = '019c0f65-58b0-7d16-83da-f360f0fa7637';
const testReleaseId = '019c0f65-58b0-7d16-83da-f360f0fa7638';
const testItemId = '019c0f65-58b0-7d16-83da-f360f0fa7639';

function testProject(): Project {
  return {
    id: testProjectId,
    name: 'Raven Game',
    description: 'A local-first game project',
    location: 'C:\\Games\\Raven',
    steam_app_id: 480,
    language: 'English',
    stage: 'concept',
    status: 'active',
    created_at: '2026-07-31T08:00:00Z',
    updated_at: '2026-07-31T08:00:00Z',
  };
}

function testTemplate(): ReleaseTemplate {
  return {
    key: 'steam-coming-soon',
    version: '1.0.0',
    schema_version: 1,
    name: 'Steam Coming Soon',
    description: 'Prepare a clear and reviewable Steam storefront.',
    reviewed_at: '2026-07-31',
    items: [
      {
        template_item_key: 'steam-app',
        title: 'Confirm the Steamworks App',
        description: 'Create or confirm the Steamworks App.',
        requirement: 'A valid App is required for the store page.',
        category: 'setup',
        requirement_level: 'required',
        source: 'platform_template',
        source_reference: 'https://partner.steamgames.com/doc/store/coming_soon',
        sort_order: 1,
      },
      {
        template_item_key: 'positioning',
        title: 'Clarify the core promise',
        description: 'Describe who the game is for.',
        requirement: 'Keep the page coherent.',
        category: 'positioning',
        requirement_level: 'recommended',
        source: 'corvus_template',
        source_reference: 'docs/product/Corvus_Studio_Launch_PRD_v0.1.md',
        sort_order: 2,
      },
    ],
  };
}

function testRelease(): ReleaseGoal {
  return {
    id: testReleaseId,
    project_id: testProjectId,
    goal_type: 'steam_coming_soon',
    title: 'Steam Coming Soon',
    status: 'preparing',
    template_key: 'steam-coming-soon',
    template_version: '1.0.0',
    checklist_summary: {
      total: 12,
      done: 1,
      blocked: 0,
      required_total: 8,
      required_done: 1,
    },
    created_at: '2026-07-31T08:00:00Z',
    updated_at: '2026-07-31T08:00:00Z',
  };
}

function testItem(): ChecklistItem {
  return {
    id: testItemId,
    release_goal_id: testReleaseId,
    title: 'Final clarity review',
    description: 'Ask someone unfamiliar with the game to review the page.',
    requirement: 'Confirm the page communicates its promise.',
    category: 'review',
    requirement_level: 'recommended',
    source: 'user',
    source_reference: '',
    status: 'in_progress',
    sort_order: 13,
    created_at: '2026-07-31T08:00:00Z',
    updated_at: '2026-07-31T08:00:00Z',
  };
}
