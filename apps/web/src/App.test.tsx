import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import App from './App';
import { projectsApi, type Project } from './api/projects';

vi.mock('./api/projects', () => ({
  projectsApi: {
    list: vi.fn(),
    create: vi.fn(),
    get: vi.fn(),
  },
}));

const mockedProjectsApi = vi.mocked(projectsApi);

describe('Project application', () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  afterEach(() => {
    cleanup();
  });

  it('shows the Project empty state', async () => {
    mockedProjectsApi.list.mockResolvedValue([]);
    renderApp('/projects');

    expect(await screen.findByText('No projects yet')).toBeInTheDocument();
    expect(screen.getAllByRole('link', { name: 'Create project' })).not.toHaveLength(0);
  });

  it('lists projects and opens a detail route', async () => {
    const project = testProject();
    mockedProjectsApi.list.mockResolvedValue([project]);
    mockedProjectsApi.get.mockResolvedValue(project);
    renderApp('/projects');

    const projectLink = await screen.findByRole('link', { name: /Raven Game/ });
    fireEvent.click(projectLink);

    expect(await screen.findByRole('heading', { name: 'Raven Game' })).toBeInTheDocument();
    expect(mockedProjectsApi.get).toHaveBeenCalledWith(project.id);
    expect(screen.getByText(project.location)).toBeInTheDocument();
  });

  it('creates a project and navigates to its detail page', async () => {
    const project = testProject();
    mockedProjectsApi.create.mockResolvedValue(project);
    renderApp('/projects/new');

    fireEvent.change(screen.getByLabelText('Game name'), { target: { value: project.name } });
    fireEvent.change(screen.getByLabelText('Project location'), {
      target: { value: project.location },
    });
    fireEvent.click(screen.getByRole('button', { name: 'Create project' }));

    await waitFor(() => {
      expect(mockedProjectsApi.create.mock.calls[0]?.[0]).toEqual(
        expect.objectContaining({
          name: project.name,
          location: project.location,
          language: 'English',
          stage: 'concept',
        }),
      );
    });
    expect(await screen.findByRole('heading', { name: project.name })).toBeInTheDocument();
    expect(mockedProjectsApi.get).not.toHaveBeenCalled();
  });

  it('reloads a direct Project detail URL from Core', async () => {
    const project = testProject();
    mockedProjectsApi.get.mockResolvedValue(project);
    renderApp(`/projects/${project.id}`);

    expect(await screen.findByRole('heading', { name: project.name })).toBeInTheDocument();
    expect(mockedProjectsApi.get).toHaveBeenCalledWith(project.id);
  });

  it('shows a retryable Core error', async () => {
    mockedProjectsApi.list.mockRejectedValue(new Error('Core is offline'));
    renderApp('/projects');

    expect(await screen.findByText('Projects are unavailable')).toBeInTheDocument();
    expect(screen.getByText('Core is offline')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Retry' })).toBeInTheDocument();
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

function testProject(): Project {
  return {
    id: '019c0f65-58b0-7d16-83da-f360f0fa7637',
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
