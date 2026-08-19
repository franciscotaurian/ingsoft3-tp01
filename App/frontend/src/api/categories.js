import { fetchApi } from './client';

export function getCategories() {
  return fetchApi('/api/categories');
}

export function createCategory(name, authHeader) {
  return fetchApi('/api/admin/categories', {
    method: 'POST',
    headers: { Authorization: authHeader },
    body: JSON.stringify({ name }),
  });
}

export function updateCategory(id, name, authHeader) {
  return fetchApi(`/api/admin/categories/${id}`, {
    method: 'PUT',
    headers: { Authorization: authHeader },
    body: JSON.stringify({ name }),
  });
}

export function deleteCategory(id, authHeader) {
  return fetchApi(`/api/admin/categories/${id}`, {
    method: 'DELETE',
    headers: { Authorization: authHeader },
  });
}
