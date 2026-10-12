import { test, expect } from '@playwright/test';

// La API tiene su propia dirección (en QA será la de Render; local es 8080)
const API = process.env.API_BASE_URL || 'http://localhost:8080';
// Auth básica para endpoints administrativos
const ADMIN_AUTH = 'Basic ' + Buffer.from('admin:admin123').toString('base64');

// Función auxiliar para listar productos disponibles desde la base real
async function listarProductos(request) {
    const res = await request.get(`${API}/api/products`);
    expect(res.status()).toBe(200);
    return await res.json();
}

test('alta de producto guarda en la BD de verdad y el borrado lo elimina', async ({ request }) => {
    const nombreUnico = `Producto Test ${Date.now()}`;

    // 1. Crear producto vía endpoint admin
    const resAlta = await request.post(`${API}/api/admin/products`, {
        headers: { Authorization: ADMIN_AUTH },
        data: {
            name: nombreUnico,
            description: 'Creado por test de integración',
            price: 1500.50,
            stock: 10,
            category_id: 1,
            image_url: ''
        }
    });
    expect(resAlta.status()).toBe(201);
    const prodCreado = await resAlta.json();
    expect(prodCreado.id).toBeDefined();

    // 2. Comprobar que la BD real lo devuelve en el listado público
    const lista = await listarProductos(request);
    const existeEnBd = lista.some(p => p.name === nombreUnico);
    expect(existeEnBd).toBe(true);

    // 3. Borrar el producto creado para limpiar la BD
    const resBorrado = await request.delete(`${API}/api/admin/products/${prodCreado.id}`, {
        headers: { Authorization: ADMIN_AUTH }
    });
    expect(resBorrado.status()).toBe(200);

    // 4. Comprobar que ya no existe en la BD
    const listaDespues = await listarProductos(request);
    const sigueExistiendo = listaDespues.some(p => p.name === nombreUnico);
    expect(sigueExistiendo).toBe(false);
});

test('producto con datos inválidos es rechazado por la API con 400 y no crea nada', async ({ request }) => {
    const antes = await listarProductos(request);

    // Intentamos crear producto con precio inválido o nombre vacío
    const res = await request.post(`${API}/api/admin/products`, {
        headers: { Authorization: ADMIN_AUTH },
        data: {
            name: '',
            price: -100,
            stock: 5,
            category_id: 1
        }
    });
    expect(res.status()).toBe(400);

    // Comprobar que la cantidad de productos en la base no cambió
    const despues = await listarProductos(request);
    expect(despues.length).toBe(antes.length);
});

test('crear pedido con teléfono inválido es rechazado por la lógica de negocio', async ({ request }) => {
    const lista = await listarProductos(request);
    test.skip(lista.length === 0, 'No hay productos disponibles para probar pedidos');

    const producto = lista[0];

    // Intentamos crear pedido con teléfono que tiene letras
    const res = await request.post(`${API}/api/orders`, {
        data: {
            customer_name: 'Cliente Test',
            customer_phone: '11-INVALIDO-TELEFONO',
            customer_address: 'Calle Falsa 123',
            items: [
                { product_id: producto.id, quantity: 1 }
            ]
        }
    });
    expect(res.status()).toBe(400);
});