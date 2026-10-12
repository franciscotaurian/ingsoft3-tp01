import { defineConfig } from '@playwright/test'

export default defineConfig({
    testDir: './e2e',                // recoge los DOS archivos; cada job correrá el suyo
    timeout: 60_000,                 // tope de CADA test (generoso para absorber cold starts)
    expect: { timeout: 15_000 },     // tope de CADA aserción
    use: {
        baseURL: process.env.E2E_BASE_URL || 'http://localhost:3000', // El FRONTEND
        trace: 'on-first-retry',         // graba la traza en el retry si falla
        screenshot: 'only-on-failure',   // captura de pantalla solo si falla
    },
    retries: 1,                      // 1 retry para tolerar demoras de red en free tier
    reporter: [['html', { open: 'never' }], ['list']],
})