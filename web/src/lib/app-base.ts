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
/**
 * Deployment base path of the web console.
 *
 * Rsbuild injects `import.meta.env.ASSET_PREFIX` from `output.assetPrefix`
 * (`/aigateway/` for the sub-path build, empty for the dev server). Every URL
 * the browser builds for console pages, backend API calls and OAuth callbacks
 * carries the same prefix, so the reverse gateway can strip it before the Go
 * server sees the request.
 */
const injectedAssetPrefix = import.meta.env.ASSET_PREFIX

export const APP_BASE_PATH = (injectedAssetPrefix ?? '').replace(/\/+$/, '')

/** Prefix a root-relative application path with the deployment base path. */
export function withBasePath(path: string): string {
  if (!APP_BASE_PATH) return path
  return `${APP_BASE_PATH}${path.startsWith('/') ? path : `/${path}`}`
}

/** Absolute base URL of the console, e.g. `https://example.com/aigateway`. */
export function getAppBaseUrl(): string {
  return `${window.location.origin}${APP_BASE_PATH}`
}
