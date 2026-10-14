// src/api/orders.test.js
// Tests unitarios para el módulo de orders.js
// Sigue estrictamente el patrón AAA (Arrange - Act - Assert)

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { cancelOrder, createOrder, updateOrderStatus } from './orders';
import * as client from './client';

describe('orders API - cancelOrder (AAA)', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('llama al endpoint POST /api/admin/orders/:id/cancel con header de autorización', async () => {
    // Arrange
    const orderId = 42;
    const authHeader = 'Basic YWRtaW46YWRtaW4xMjM=';
    const mockResponse = { id: 42, status: 'cancelado' };
    const fetchApiSpy = vi.spyOn(client, 'fetchApi').mockResolvedValue(mockResponse);

    // Act
    const result = await cancelOrder(orderId, authHeader);

    // Assert
    expect(fetchApiSpy).toHaveBeenCalledTimes(1);
    expect(fetchApiSpy).toHaveBeenCalledWith(`/api/admin/orders/${orderId}/cancel`, {
      method: 'POST',
      headers: { Authorization: authHeader },
    });
    expect(result).toEqual(mockResponse);
  });

  it('propaga el error si el servidor rechaza la cancelación', async () => {
    // Arrange
    const orderId = 42;
    const authHeader = 'Basic YWRtaW46YWRtaW4xMjM=';
    vi.spyOn(client, 'fetchApi').mockRejectedValue(new Error('no se puede cancelar un pedido en estado entregado'));

    // Act & Assert
    await expect(cancelOrder(orderId, authHeader))
      .rejects
      .toThrow('no se puede cancelar un pedido en estado entregado');
  });
});
