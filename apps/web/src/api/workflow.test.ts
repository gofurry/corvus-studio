import { afterEach, describe, expect, it, vi } from 'vitest';

import { checklistApi, type ChecklistItem } from './checklist';
import { releasesApi, type ReleaseTemplate } from './releases';

describe('release and checklist APIs', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('loads a versioned release template through the generated client', async () => {
    const template = testTemplate();
    const fetchMock = vi.fn().mockResolvedValue(
      new Response(JSON.stringify(template), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }),
    );
    vi.stubGlobal('fetch', fetchMock);

    await expect(releasesApi.getTemplate(template.key)).resolves.toEqual(template);
    const request = fetchMock.mock.calls[0]?.[0] as Request;
    expect(request.method).toBe('GET');
    expect(request.url).toContain('/api/v1/release-templates/steam-coming-soon');
  });

  it('sends checklist filters and returns generated client data', async () => {
    const item = testChecklistItem();
    const fetchMock = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ items: [item] }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }),
    );
    vi.stubGlobal('fetch', fetchMock);

    await expect(
      checklistApi.list(item.release_goal_id, {
        status: 'in_progress',
        source: 'user',
        category: 'review',
      }),
    ).resolves.toEqual([item]);

    const request = fetchMock.mock.calls[0]?.[0] as Request;
    expect(request.url).toContain('status=in_progress');
    expect(request.url).toContain('source=user');
    expect(request.url).toContain('category=review');
  });

  it('normalizes release workflow error envelopes', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(
        new Response(
          JSON.stringify({
            error: {
              code: 'release_not_ready',
              message: 'Required tasks are not complete',
              recoverable: true,
            },
          }),
          { status: 409, headers: { 'Content-Type': 'application/json' } },
        ),
      ),
    );

    await expect(releasesApi.transition(testReleaseId, 'ready_for_review')).rejects.toMatchObject({
      name: 'WorkflowApiError',
      code: 'release_not_ready',
      message: 'Required tasks are not complete',
      recoverable: true,
    });
  });
});

const testReleaseId = '019c0f65-58b0-7d16-83da-f360f0fa7638';

function testTemplate(): ReleaseTemplate {
  return {
    key: 'steam-coming-soon',
    version: '1.0.0',
    schema_version: 1,
    name: 'Steam Coming Soon',
    description: 'Prepare a clear Coming Soon page.',
    reviewed_at: '2026-07-31',
    items: [],
  };
}

function testChecklistItem(): ChecklistItem {
  return {
    id: '019c0f65-58b0-7d16-83da-f360f0fa7639',
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
