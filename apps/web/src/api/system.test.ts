import { afterEach, describe, expect, it, vi } from 'vitest';

import { systemApi } from './system';

describe('systemApi', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('returns the selected absolute directory', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ selected: true, path: 'C:\\Games\\Raven' }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }),
    );
    vi.stubGlobal('fetch', fetchMock);

    await expect(systemApi.selectProjectDirectory()).resolves.toBe('C:\\Games\\Raven');
    expect(fetchMock).toHaveBeenCalledOnce();
    const request = fetchMock.mock.calls[0]?.[0] as Request;
    await expect(request.clone().json()).resolves.toEqual({ purpose: 'project_location' });
  });

  it('returns null when the user cancels selection', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(
        new Response(JSON.stringify({ selected: false, path: null }), {
          status: 200,
          headers: { 'Content-Type': 'application/json' },
        }),
      ),
    );

    await expect(systemApi.selectProjectDirectory()).resolves.toBeNull();
  });

  it('normalizes an unavailable picker error', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(
        new Response(
          JSON.stringify({
            error: {
              code: 'directory_picker_unavailable',
              message: 'Native directory picker is unavailable',
              recoverable: true,
            },
          }),
          {
            status: 503,
            headers: { 'Content-Type': 'application/json' },
          },
        ),
      ),
    );

    await expect(systemApi.selectProjectDirectory()).rejects.toMatchObject({
      name: 'DirectoryPickerApiError',
      code: 'directory_picker_unavailable',
      recoverable: true,
    });
  });
});
