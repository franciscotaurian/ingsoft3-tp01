// src/lib/pedido.test.js
// Suite de tests para la lógica pura de pedidos.
//
// Técnicas del TP:
//   ✅ Parametrizado : 'rechaza cuando %s' con it.each (Tests F-1)
//   ✅ Caso de error  : borde exacto del nombre mínimo (Test F-2)

import { describe, it, expect } from 'vitest'
import {
  validarFormularioPedido,
  calcularTotalPedido,
  MIN_NOMBRE_LENGTH,
} from './pedido.js'

// ─────────────────────────────────────────────────────────────────────────────
// Test F-1 — Parametrizado: rechaza datos inválidos del formulario
// ─────────────────────────────────────────────────────────────────────────────

describe('validarFormularioPedido', () => {
  it.each([
    ['nombre vacío',        { nombre: '',          telefono: '12345', direccion: 'Av. X', items: [1] }],
    ['nombre muy corto',    { nombre: 'AB',         telefono: '12345', direccion: 'Av. X', items: [1] }],
    ['teléfono con guion',  { nombre: 'Juan Pérez', telefono: '123-4', direccion: 'Av. X', items: [1] }],
    ['teléfono con letras', { nombre: 'Juan Pérez', telefono: 'abc12', direccion: 'Av. X', items: [1] }],
    ['dirección vacía',     { nombre: 'Juan Pérez', telefono: '12345', direccion: '',      items: [1] }],
    ['sin items',           { nombre: 'Juan Pérez', telefono: '12345', direccion: 'Av. X', items: [] }],
  ])('rechaza cuando %s', (_caso, formulario) => {
    // Act
    const resultado = validarFormularioPedido(formulario)

    // Assert
    expect(resultado.valido).toBe(false)
  })

  // ─────────────────────────────────────────────────────────────────────────
  // Test F-2 — Caso de error: borde exacto del largo mínimo de nombre
  //
  // Si alguien cambia `< MIN_NOMBRE_LENGTH` por `<= MIN_NOMBRE_LENGTH`,
  // el segundo expect falla porque el nombre exacto quedaría rechazado.
  // ─────────────────────────────────────────────────────────────────────────

  it('rechaza nombre con largo menor al mínimo y acepta con el largo exacto', () => {
    // Arrange
    const nombreExacto = 'A'.repeat(MIN_NOMBRE_LENGTH)       // borde: debe pasar
    const nombreCorto  = 'A'.repeat(MIN_NOMBRE_LENGTH - 1)   // borde - 1: debe fallar
    const base = { telefono: '12345', direccion: 'Av. X', items: [1] }

    // Act + Assert: un carácter menos es inválido
    expect(validarFormularioPedido({ nombre: nombreCorto, ...base }).valido).toBe(false)

    // Act + Assert: el largo exacto es válido
    expect(validarFormularioPedido({ nombre: nombreExacto, ...base }).valido).toBe(true)
  })

  it('acepta un formulario con todos los datos correctos', () => {
    // Arrange
    const formulario = {
      nombre:    'Juan Pérez',
      telefono:  '3516789012',
      direccion: 'Av. Siempreviva 742',
      items:     [{ id: 1, cantidad: 2 }],
    }

    // Act
    const resultado = validarFormularioPedido(formulario)

    // Assert
    expect(resultado.valido).toBe(true)
    expect(resultado.errores).toEqual({})
  })
})

// ─────────────────────────────────────────────────────────────────────────────
// Test F-3 — Lógica pura: calcularTotalPedido
// ─────────────────────────────────────────────────────────────────────────────

describe('calcularTotalPedido', () => {
  it('devuelve 0 si la lista está vacía o es nula', () => {
    // Arrange + Act + Assert
    expect(calcularTotalPedido([])).toBe(0)
    expect(calcularTotalPedido(null)).toBe(0)
  })

  it('suma precio × cantidad de cada ítem correctamente', () => {
    // Arrange
    const items = [
      { precio: 500, cantidad: 2 },  // 1000
      { precio: 300, cantidad: 3 },  // 900
    ]

    // Act
    const total = calcularTotalPedido(items)

    // Assert
    expect(total).toBe(1900)
  })
})
