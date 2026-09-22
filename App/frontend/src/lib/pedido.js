// src/lib/pedido.js
/** Longitud mínima de caracteres para el nombre del cliente. */
export const MIN_NOMBRE_LENGTH = 3

/**
 * Valida los datos del formulario de pedido antes de enviarlo a la API.
 *
 * @param {{ nombre: string, telefono: string, direccion: string, items: Array }} param0
 * @returns {{ valido: boolean, errores: Record<string, string> }}
 */
export function validarFormularioPedido({ nombre, telefono, direccion, items }) {
  const errores = {}

  // R-A: nombre obligatorio con largo mínimo
  if (!nombre || nombre.trim().length < MIN_NOMBRE_LENGTH) {
    errores.nombre = `El nombre debe tener al menos ${MIN_NOMBRE_LENGTH} caracteres.`
  }

  // R-B: teléfono solo dígitos (misma regla que el backend)
  if (!telefono || !/^\d+$/.test(telefono.trim())) {
    errores.telefono = 'El teléfono debe contener solo números.'
  }

  // R-A: dirección obligatoria
  if (!direccion || direccion.trim() === '') {
    errores.direccion = 'La dirección es obligatoria.'
  }

  // R-C: al menos un ítem
  if (!items || items.length === 0) {
    errores.items = 'Debe agregar al menos un producto al pedido.'
  }

  return { valido: Object.keys(errores).length === 0, errores }
}

/**
 * Calcula el total del pedido sumando precio × cantidad de cada ítem.
 *
 * @param {Array<{ precio: number, cantidad: number }>} items
 * @returns {number}
 */
export function calcularTotalPedido(items) {
  if (!items || items.length === 0) return 0
  return items.reduce((acc, item) => acc + item.precio * item.cantidad, 0)
}
