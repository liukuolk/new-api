/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { afterEach, describe, expect, it, vi } from 'vitest'

describe('app-base', () => {
  afterEach(() => {
    vi.unstubAllEnvs()
    vi.resetModules()
  })

  it('keeps root-relative paths when no deployment prefix is built in', async () => {
    const { APP_BASE_PATH, withBasePath, getAppBaseUrl } = await import(
      '../app-base'
    )

    expect(APP_BASE_PATH).toBe('')
    expect(withBasePath('/api/status')).toBe('/api/status')
    expect(getAppBaseUrl()).toBe(window.location.origin)
  })

  it('prefixes console, API and OAuth URLs with the deployment base path', async () => {
    vi.stubEnv('ASSET_PREFIX', '/aigateway/')

    const { APP_BASE_PATH, withBasePath, getAppBaseUrl } = await import(
      '../app-base'
    )

    expect(APP_BASE_PATH).toBe('/aigateway')
    expect(withBasePath('/api/status')).toBe('/aigateway/api/status')
    expect(withBasePath('oauth/github')).toBe('/aigateway/oauth/github')
    expect(getAppBaseUrl()).toBe(`${window.location.origin}/aigateway`)
  })
})
