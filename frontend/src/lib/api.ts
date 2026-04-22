const apiBaseUrl = (import.meta.env.VITE_API_URL ?? "").replace(/\/+$/, "");

export function apiFetch(path: string, init?: RequestInit) {
  return fetch(buildApiUrl(path), init);
}

function buildApiUrl(path: string) {
  if (path.startsWith("http://") || path.startsWith("https://")) {
    return path;
  }

  if (!path.startsWith("/")) {
    return `${apiBaseUrl}/${path}`;
  }

  return `${apiBaseUrl}${path}`;
}
