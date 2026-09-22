package service

// order_service_test.go — Suite de tests para OrderService.
//
// Cubre las siguientes reglas de negocio:
//   R-A: Datos del cliente obligatorios (nombre, teléfono, dirección)
//   R-B: Formato del teléfono (solo dígitos)
//   R-C: El pedido debe tener al menos un ítem
//   R-E: No se puede superar el stock disponible
//   R-F: Transiciones de estado permitidas del pedido
//
// Técnicas del TP:
//   ✅ Parametrizado : TestCreate_ValidaQueLosDatosDelClienteSeanObligatorios,
//                      TestCreate_ValidaFormatoTelefono,
//                      TestUpdateStatus_TransicionesDeEstadoValidas,
//                      TestUpdateStatus_TransicionInvalidaEsRechazada
//   ✅ Caso de error  : TestCreate_ValidaFormatoTelefono, TestCreate_PedidoSinItemsEsRechazado,
//                      TestUpdateStatus_TransicionInvalidaEsRechazada
//   ✅ Mock           : TestCreate_LlamaAlRepositorioConStockDescontado

import (
	"testing"

	"realico-comidas-backend/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ─────────────────────────────────────────────────────────────────────────────
// Test 1 — Parametrizado: datos del cliente obligatorios (R-A)
// ─────────────────────────────────────────────────────────────────────────────

func TestCreate_ValidaQueLosDatosDelClienteSeanObligatorios(t *testing.T) {
	casos := []struct {
		nombre        string
		dto           CreateOrderDTO
		errorEsperado error
	}{
		{
			nombre: "nombre_vacio",
			dto: CreateOrderDTO{
				CustomerName:    "",
				CustomerPhone:   "1234567890",
				CustomerAddress: "Av. Siempreviva 742",
				Items:           []CreateOrderItemDTO{{ProductID: 1, Quantity: 1}},
			},
			errorEsperado: ErrCustomerDataRequired,
		},
		{
			nombre: "nombre_solo_espacios",
			dto: CreateOrderDTO{
				CustomerName:    "   ",
				CustomerPhone:   "1234567890",
				CustomerAddress: "Av. Siempreviva 742",
				Items:           []CreateOrderItemDTO{{ProductID: 1, Quantity: 1}},
			},
			errorEsperado: ErrCustomerDataRequired,
		},
		{
			nombre: "telefono_vacio",
			dto: CreateOrderDTO{
				CustomerName:    "Juan Pérez",
				CustomerPhone:   "",
				CustomerAddress: "Av. Siempreviva 742",
				Items:           []CreateOrderItemDTO{{ProductID: 1, Quantity: 1}},
			},
			errorEsperado: ErrCustomerDataRequired,
		},
		{
			nombre: "direccion_vacia",
			dto: CreateOrderDTO{
				CustomerName:    "Juan Pérez",
				CustomerPhone:   "1234567890",
				CustomerAddress: "",
				Items:           []CreateOrderItemDTO{{ProductID: 1, Quantity: 1}},
			},
			errorEsperado: ErrCustomerDataRequired,
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// Arrange — los repos no se invocan: la validación falla antes de usarlos
			svc := NewOrderService(&MockOrderRepository{}, &MockProductRepository{})

			// Act
			_, err := svc.Create(c.dto)

			// Assert
			assert.ErrorIs(t, err, c.errorEsperado)
		})
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Test 2 — Caso de error: formato de teléfono inválido (R-B)
// ─────────────────────────────────────────────────────────────────────────────

func TestCreate_ValidaFormatoTelefono(t *testing.T) {
	casos := []struct {
		nombre   string
		telefono string
	}{
		{"con_guion", "123-456-7890"},
		{"con_espacio", "123 456 7890"},
		{"con_letra", "123ABC7890"},
		{"con_parentesis", "(123)4567890"},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// Arrange
			svc := NewOrderService(&MockOrderRepository{}, &MockProductRepository{})
			dto := CreateOrderDTO{
				CustomerName:    "Juan Pérez",
				CustomerPhone:   c.telefono,
				CustomerAddress: "Av. Siempreviva 742",
				Items:           []CreateOrderItemDTO{{ProductID: 1, Quantity: 1}},
			}

			// Act
			_, err := svc.Create(dto)

			// Assert
			assert.ErrorIs(t, err, ErrInvalidPhone)
		})
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Test 3 — Caso de error: pedido sin ítems (R-C)
// ─────────────────────────────────────────────────────────────────────────────

func TestCreate_PedidoSinItemsEsRechazado(t *testing.T) {
	// Arrange
	svc := NewOrderService(&MockOrderRepository{}, &MockProductRepository{})
	dto := CreateOrderDTO{
		CustomerName:    "Juan Pérez",
		CustomerPhone:   "1234567890",
		CustomerAddress: "Av. Siempreviva 742",
		Items:           []CreateOrderItemDTO{}, // lista vacía
	}

	// Act
	_, err := svc.Create(dto)

	// Assert
	assert.ErrorIs(t, err, ErrOrderNoItems)
}

// ─────────────────────────────────────────────────────────────────────────────
// Test 4 — Caso de error con mock: stock insuficiente (R-E)
// ─────────────────────────────────────────────────────────────────────────────

func TestCreate_FallaSiStockInsuficiente(t *testing.T) {
	// Arrange
	mockOrderRepo := &MockOrderRepository{}
	mockProductRepo := &MockProductRepository{}

	// El mock devuelve un producto con stock=2
	productoConStock2 := &models.Product{
		ID:    1,
		Name:  "Empanada",
		Price: 500,
		Stock: 2,
	}
	mockProductRepo.On("FindByID", uint(1)).Return(productoConStock2, nil)

	svc := NewOrderService(mockOrderRepo, mockProductRepo)
	dto := CreateOrderDTO{
		CustomerName:    "Juan Pérez",
		CustomerPhone:   "1234567890",
		CustomerAddress: "Av. Siempreviva 742",
		Items:           []CreateOrderItemDTO{{ProductID: 1, Quantity: 5}}, // pide 5, hay 2
	}

	// Act
	_, err := svc.Create(dto)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "stock")
	mockProductRepo.AssertExpectations(t)
}

// ─────────────────────────────────────────────────────────────────────────────
// Test 5 — Mock (verifica interacción): stock descontado correctamente (R-E + R-3)
//
// Este test NO verifica el valor de retorno — verifica que CreateWithTx fue
// llamado con el mapa de stock correcto {productID: stockRestante}.
// Es el equivalente del Verify(n => n.Enviar(...)) de Moq.
// ─────────────────────────────────────────────────────────────────────────────

func TestCreate_LlamaAlRepositorioConStockDescontado(t *testing.T) {
	// Arrange
	mockOrderRepo := &MockOrderRepository{}
	mockProductRepo := &MockProductRepository{}

	producto := &models.Product{ID: 1, Name: "Empanada", Price: 500.0, Stock: 10}
	mockProductRepo.On("FindByID", uint(1)).Return(producto, nil)

	// Verificamos que CreateWithTx recibe exactamente el mapa con stock=7 (10-3)
	stockEsperado := map[uint]int{uint(1): 7}
	mockOrderRepo.
		On("CreateWithTx", mock.AnythingOfType("*models.Order"), stockEsperado).
		Return(nil)
	mockOrderRepo.
		On("FindByID", mock.AnythingOfType("uint")).
		Return(&models.Order{ID: 1, Status: models.OrderStatusPendiente}, nil)

	svc := NewOrderService(mockOrderRepo, mockProductRepo)
	dto := CreateOrderDTO{
		CustomerName:    "Juan Pérez",
		CustomerPhone:   "1234567890",
		CustomerAddress: "Av. Siempreviva 742",
		Items:           []CreateOrderItemDTO{{ProductID: 1, Quantity: 3}},
	}

	// Act
	_, err := svc.Create(dto)

	// Assert
	assert.NoError(t, err)
	// Si la fórmula de descuento de stock cambia, AssertExpectations falla
	// porque el mapa no coincide con {1: 7}
	mockOrderRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

// ─────────────────────────────────────────────────────────────────────────────
// Test 6 — Parametrizado: transiciones de estado válidas (R-F)
// ─────────────────────────────────────────────────────────────────────────────

func TestUpdateStatus_TransicionesDeEstadoValidas(t *testing.T) {
	casos := []struct {
		nombre       string
		estadoActual string
		nuevoEstado  string
	}{
		{"pendiente_a_confirmado", models.OrderStatusPendiente, models.OrderStatusConfirmado},
		{"confirmado_a_entregado", models.OrderStatusConfirmado, models.OrderStatusEntregado},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// Arrange
			mockOrderRepo := &MockOrderRepository{}
			order := &models.Order{ID: 1, Status: c.estadoActual}

			// FindByID se llama 2 veces: antes de actualizar y al retornar el resultado
			mockOrderRepo.On("FindByID", uint(1)).Return(order, nil)
			mockOrderRepo.On("UpdateStatus", uint(1), c.nuevoEstado).Return(nil)

			svc := NewOrderService(mockOrderRepo, &MockProductRepository{})

			// Act
			resultado, err := svc.UpdateStatus(1, c.nuevoEstado)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, resultado)
		})
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Test 7 — Parametrizado + caso de error: transiciones inválidas (R-F)
// ─────────────────────────────────────────────────────────────────────────────

func TestUpdateStatus_TransicionInvalidaEsRechazada(t *testing.T) {
	casos := []struct {
		nombre       string
		estadoActual string
		nuevoEstado  string
	}{
		{"entregado_a_pendiente", models.OrderStatusEntregado, models.OrderStatusPendiente},
		{"pendiente_a_entregado", models.OrderStatusPendiente, models.OrderStatusEntregado},
		{"confirmado_a_pendiente", models.OrderStatusConfirmado, models.OrderStatusPendiente},
		{"entregado_a_confirmado", models.OrderStatusEntregado, models.OrderStatusConfirmado},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// Arrange
			mockOrderRepo := &MockOrderRepository{}
			order := &models.Order{ID: 1, Status: c.estadoActual}
			mockOrderRepo.On("FindByID", uint(1)).Return(order, nil)

			svc := NewOrderService(mockOrderRepo, &MockProductRepository{})

			// Act
			_, err := svc.UpdateStatus(1, c.nuevoEstado)

			// Assert — si alguien habilita una transición inválida, este test se pone rojo
			assert.ErrorIs(t, err, ErrInvalidStatusChange)
		})
	}
}
