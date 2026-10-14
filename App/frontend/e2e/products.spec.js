import { test, expect } from '@playwright/test';

test('el usuario puede ver productos, agregarlos al carrito y realizar un pedido (con cleanup admin de restauración de stock)', async ({ page }) => {
    // 1. Navegar a la página principal
    await page.goto('/');

    // Esperar a que cargue el título principal
    await expect(page.getByRole('heading', { name: 'Realico Comidas' })).toBeVisible();

    // 2. Registrar el stock inicial del primer producto para verificar su inmutabilidad al final
    const primeraCard = page.locator('.product-card').first();
    await expect(primeraCard).toBeVisible();
    const stockInicialTexto = await primeraCard.locator('.product-stock').innerText();

    // 3. Buscar el botón "Agregar al Carrito" del producto y clickearlo
    const btnAgregar = primeraCard.getByRole('button', { name: 'Agregar al Carrito' });
    await expect(btnAgregar).toBeVisible();
    await btnAgregar.click();

    // 4. Verificar que el carrito muestre el formulario y los datos de entrega
    await expect(page.getByRole('heading', { name: 'Tu Pedido' })).toBeVisible();
    await expect(page.getByRole('heading', { name: 'Datos de Entrega' })).toBeVisible();

    // 5. Llenar los datos del cliente con los placeholders reales
    const nombreCliente = `Usuario e2e ${Date.now()}`;
    await page.getByPlaceholder('Ej: Juan Pérez').fill(nombreCliente);
    await page.getByPlaceholder('Ej: 1122334455').fill('2302123456');
    await page.getByPlaceholder('Ej: Av. Siempreviva 123').fill('Av. Principal 789');

    // 6. Confirmar y enviar el pedido (botón real: "Hacer Pedido por WhatsApp")
    const btnEnviar = page.getByRole('button', { name: 'Hacer Pedido por WhatsApp' });
    await expect(btnEnviar).toBeVisible();
    await btnEnviar.click();

    // 7. Verificar el mensaje de éxito en pantalla y que el carrito se vació
    await expect(page.getByText(/¡Pedido confirmado correctamente!/i)).toBeVisible();
    await expect(page.getByText('El carrito está vacío')).toBeVisible();

    // 8. Cleanup (dejar la base como se encontró): El admin cancela el pedido restaurando el stock
    await page.getByRole('button', { name: 'Panel Admin' }).click();
    await expect(page.getByRole('heading', { name: 'Acceso de Administrador' })).toBeVisible();

    await page.getByPlaceholder('admin').fill('admin');
    await page.getByPlaceholder('••••••••').fill('admin123');
    await page.getByRole('button', { name: 'Ingresar' }).click();

    await expect(page.getByRole('heading', { name: 'Panel de Administración' })).toBeVisible();

    // Localizar la fila del pedido recién creado
    const filaPedido = page.locator('tr', { hasText: nombreCliente });
    await expect(filaPedido).toBeVisible();

    // Aceptar el diálogo de confirmación window.confirm antes del click
    page.once('dialog', dialog => dialog.accept());

    // Click en cancelar pedido para restaurar stock
    await filaPedido.getByRole('button', { name: /Cancelar/i }).click();

    // Verificar que el estado del pedido pase a cancelado (apuntando unívocamente al badge de estado)
    await expect(filaPedido.locator('.badge-cancelado')).toBeVisible();

    // Cerrar sesión para dejar limpio el estado del navegador
    await page.getByRole('button', { name: 'Cerrar Sesión' }).click();

    // 9. Aserción de Inmutabilidad de Stock: Regresar al catálogo y verificar que el stock volvió exactamente a su valor inicial
    await page.getByRole('button', { name: 'Catálogo' }).click();
    await expect(primeraCard).toBeVisible();
    await expect(primeraCard.locator('.product-stock')).toHaveText(stockInicialTexto);
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

