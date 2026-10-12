import { test, expect } from '@playwright/test';

test('el usuario puede ver productos, agregarlos al carrito y realizar un pedido', async ({ page }) => {
    // 1. Navegar a la página principal
    await page.goto('/');

    // Esperar a que cargue el título principal
    await expect(page.getByRole('heading', { name: 'Realico Comidas' })).toBeVisible();

    // 2. Buscar el primer botón "Agregar al Carrito" y clickearlo
    const btnAgregar = page.getByRole('button', { name: 'Agregar al Carrito' }).first();
    await expect(btnAgregar).toBeVisible();
    await btnAgregar.click();

    // 3. Verificar que el carrito muestre el formulario y los datos de entrega
    await expect(page.getByRole('heading', { name: 'Tu Pedido' })).toBeVisible();
    await expect(page.getByRole('heading', { name: 'Datos de Entrega' })).toBeVisible();

    // 4. Llenar los datos del cliente con los placeholders reales
    const nombreCliente = `Usuario e2e ${Date.now()}`;
    await page.getByPlaceholder('Ej: Juan Pérez').fill(nombreCliente);
    await page.getByPlaceholder('Ej: 1122334455').fill('2302123456');
    await page.getByPlaceholder('Ej: Av. Siempreviva 123').fill('Av. Principal 789');

    // 5. Confirmar y enviar el pedido (botón real: "Hacer Pedido por WhatsApp")
    const btnEnviar = page.getByRole('button', { name: 'Hacer Pedido por WhatsApp' });
    await expect(btnEnviar).toBeVisible();
    await btnEnviar.click();

    // 6. Verificar el mensaje de éxito en pantalla y que el carrito se vació
    await expect(page.getByText(/¡Pedido confirmado correctamente!/i)).toBeVisible();
    await expect(page.getByText('El carrito está vacío')).toBeVisible();
});

test('validación del formulario de pedido: intentar enviar sin datos requeridos muestra alerta', async ({ page }) => {
    await page.goto('/');

    // Agregar un producto al carrito
    const btnAgregar = page.getByRole('button', { name: 'Agregar al Carrito' }).first();
    await expect(btnAgregar).toBeVisible();
    await btnAgregar.click();

    // Dejar campos vacíos e intentar enviar
    const btnEnviar = page.getByRole('button', { name: 'Hacer Pedido por WhatsApp' });
    await expect(btnEnviar).toBeVisible();
    await btnEnviar.click();

    // El input requerido no completado queda enfocado por validación HTML5
    await expect(page.getByPlaceholder('Ej: Juan Pérez')).toBeFocused();
});

test('panel admin: login con credenciales correctas muestra el panel y permite logout', async ({ page }) => {
    await page.goto('/');

    // 1. Ir a la pestaña Panel Admin
    await page.getByRole('button', { name: 'Panel Admin' }).click();
    await expect(page.getByRole('heading', { name: 'Acceso de Administrador' })).toBeVisible();

    // 2. Ingresar credenciales de administrador (admin / admin123)
    await page.getByPlaceholder('admin').fill('admin');
    await page.getByPlaceholder('••••••••').fill('admin123');
    
    // 3. Click en el botón real: "Ingresar"
    await page.getByRole('button', { name: 'Ingresar' }).click();

    // 4. Verificar que se abrió el panel admin autenticado
    await expect(page.getByRole('heading', { name: 'Panel de Administración' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Cerrar Sesión' })).toBeVisible();

    // 5. Cerrar sesión
    await page.getByRole('button', { name: 'Cerrar Sesión' }).click();
    await expect(page.getByRole('heading', { name: 'Acceso de Administrador' })).toBeVisible();
});

