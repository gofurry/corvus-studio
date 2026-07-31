import { afterEach, describe, expect, it, vi } from 'vitest';

import { projectsApi, type Project } from './projects';

describe('projectsApi', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('returns generated client data', async () => {
    const project = testProject();
    const fetchMock = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ items: [project] }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }),
    );
    vi.stubGlobal('fetch', fetchMock);

    await expect(projectsApi.list()).resolves.toEqual([project]);
    expect(fetchMock).toHaveBeenCalledOnce();
  });

  it('normalizes the API error envelope', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      new Response(
        JSON.stringify({
          error: {
            code: 'project_location_conflict',
            message: 'A Project already references this location',
            recoverable: true,
          },
        }),
        {
          status: 409,
          headers: { 'Content-Type': 'application/json' },
        },
      ),
    );
    vi.stubGlobal('fetch', fetchMock);

    await expect(
      projectsApi.create({
        name: 'Raven',
        location: 'C:\\Games\\Raven',
        language: 'English',
        stage: 'concept',
      }),
    ).rejects.toMatchObject(
      expect.objectContaining({
        name: 'ProjectApiError',
        code: 'project_location_conflict',
        message: 'A Project already references this location',
        recoverable: true,
      }),
    );
  });
});

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
