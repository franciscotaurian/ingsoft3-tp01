import { fetchApi } from './client';

// El historial de pedidos requiere autenticación de administrador
export function getOrders(authHeader) {
  return fetchApi('/api/admin/orders', {
    headers: { Authorization: authHeader },
  });
}

export function createOrder(orderData) {
  return fetchApi('/api/orders', {
    method: 'POST',
    body: JSON.stringify(orderData),
  });
}

export function updateOrderStatus(id, status, authHeader) {
  return fetchApi(`/api/admin/orders/${id}/status`, {
    method: 'PATCH',
    headers: { Authorization: authHeader },
    body: JSON.stringify({ status }),
  });
}

export function getAdminMetrics(authHeader) {
  return fetchApi('/api/admin/metrics', {
    headers: { Authorization: authHeader },
  });
}
