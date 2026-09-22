// src/api/client.test.js
// Tests para fetchApi usando vi.fn() como mock de fetch global.
//
// Técnica del TP:
//   ✅ Mock: vi.fn() reemplaza fetch global; toHaveBeenCalledWith verifica
//            la interacción (qué URL pidió), no solo el valor devuelto.
//   ✅ Caso de error: respuesta con ok=false lanza Error con el mensaje del servidor.

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { fetchApi } from './client.js'

describe('fetchApi', () => {
  beforeEach(() => {
    // Reemplazar fetch global por un impostor — ningún test toca la red
    global.fetch = vi.fn()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  // ─────────────────────────────────────────────────────────────────────────
  // Test F-4a — Mock + caso de error: el servidor responde con error JSON
  //
  // Verifica la lógica de fetchApi: parsear el campo `error` del body
  // y lanzarlo como un Error. Si alguien cambia `data.error` por `data.message`,
  // este test se pone en rojo.
  // ─────────────────────────────────────────────────────────────────────────

  it('lanza el mensaje de error del servidor cuando la respuesta no es ok', async () => {
    // Arrange: mock simula una respuesta 400 con cuerpo JSON de error
    global.fetch.mockResolvedValue({
      ok: false,
      status: 400,
      json: vi.fn().mockResolvedValue({ error: 'el teléfono debe contener solo números' }),
    })

    // Act + Assert
    await expect(fetchApi('/api/orders', { method: 'POST' }))
      .rejects
      .toThrow('el teléfono debe contener solo números')
  })

  // ─────────────────────────────────────────────────────────────────────────
  // Test F-4b — Mock (verifica interacción): 204 No Content retorna null
  //
  // El assert de toHaveBeenCalledWith es el equivalente del Verify de Moq:
  // no verifica el valor devuelto, verifica QUÉ URL se le pidió al fetch.
  // ─────────────────────────────────────────────────────────────────────────

  it('retorna null para respuestas 204 No Content y verifica la URL llamada', async () => {
    // Arrange
    global.fetch.mockResolvedValue({ ok: true, status: 204, json: vi.fn() })

    // Act
    const resultado = await fetchApi('/api/products')

    // Assert — verifica el VALOR retornado
    expect(resultado).toBeNull()

    // Assert — verifica la INTERACCIÓN: fetch fue llamado con la URL correcta
    expect(global.fetch).toHaveBeenCalledWith(
      expect.stringContaining('/api/products'),
      expect.any(Object),
    )
  })

  // ─────────────────────────────────────────────────────────────────────────
  // Test F-4c — Mock + caso de error: body no es JSON parseable
  //
  // Cuando el servidor responde con error pero el body no es JSON,
  // fetchApi debe lanzar el mensaje genérico 'Error en la petición'.
  // ─────────────────────────────────────────────────────────────────────────

  it('lanza un error genérico si la respuesta de error no tiene JSON válido', async () => {
    // Arrange
    global.fetch.mockResolvedValue({
      ok: false,
      status: 500,
      json: vi.fn().mockRejectedValue(new SyntaxError('Unexpected token')),
    })

    // Act + Assert
    await expect(fetchApi('/api/products'))
      .rejects
      .toThrow('Error en la petición')
  })
})
