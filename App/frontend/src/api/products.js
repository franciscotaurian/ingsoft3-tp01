import { fetchApi } from './client';

export function getPublicProducts() {
  return fetchApi('/api/products');
}

export function getAdminProducts(authHeader) {
  return fetchApi('/api/admin/products', {
    headers: { Authorization: authHeader },
  });
}

export function createProduct(productData, authHeader) {
  return fetchApi('/api/admin/products', {
    method: 'POST',
    headers: { Authorization: authHeader },
    body: JSON.stringify(productData),
  });
}

export function updateProduct(id, productData, authHeader) {
  return fetchApi(`/api/admin/products/${id}`, {
    method: 'PUT',
    headers: { Authorization: authHeader },
    body: JSON.stringify(productData),
  });
}

export function deleteProduct(id, authHeader) {
  return fetchApi(`/api/admin/products/${id}`, {
    method: 'DELETE',
    headers: { Authorization: authHeader },
  });
}
