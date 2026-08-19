import React, { useState } from 'react';

export default function Cart({ cartItems, onUpdateQuantity, onRemoveItem, onSubmitOrder, loading }) {
  const [customerName, setCustomerName] = useState('');
  const [customerPhone, setCustomerPhone] = useState('');
  const [customerAddress, setCustomerAddress] = useState('');
  const [error, setError] = useState('');

  const total = cartItems.reduce((sum, item) => sum + item.price * item.quantity, 0);

  const handleSubmit = (e) => {
    e.preventDefault();
    setError('');

    if (cartItems.length === 0) {
      setError('El carrito está vacío');
      return;
    }

    const nameClean = customerName.trim();
    const phoneClean = customerPhone.trim();
    const addressClean = customerAddress.trim();

    if (!nameClean || !phoneClean || !addressClean) {
      setError('Por favor completa todos los datos del cliente');
      return;
    }

    // Validación: el teléfono debe contener solo dígitos numéricos
    if (!/^\d+$/.test(phoneClean)) {
      setError('El teléfono debe contener solo números (sin espacios, guiones ni caracteres especiales)');
      return;
    }

    onSubmitOrder({
      customer_name: nameClean,
      customer_phone: phoneClean,
      customer_address: addressClean,
      items: cartItems.map((item) => ({
        product_id: item.id,
        quantity: item.quantity,
      })),
    });
  };

  return (
    <div className="cart-sidebar">
      <h2>Tu Pedido</h2>

      {error && <div className="alert-error" style={{ marginTop: '0.5rem' }}>{error}</div>}

      {cartItems.length === 0 ? (
        <p style={{ margin: '1rem 0', color: '#777' }}>El carrito está vacío</p>
      ) : (
        <>
          <ul className="cart-items">
            {cartItems.map((item) => (
              <li key={item.id} className="cart-item">
                <div>
                  <strong>{item.name}</strong>
                  <br />
                  <small>${item.price.toFixed(2)} c/u</small>
                </div>
                <div className="cart-item-actions">
                  <button onClick={() => onUpdateQuantity(item.id, item.quantity - 1)}>-</button>
                  <span>{item.quantity}</span>
                  <button
                    onClick={() => onUpdateQuantity(item.id, item.quantity + 1)}
                    disabled={item.quantity >= item.stock}
                  >
                    +
                  </button>
                  <button
                    onClick={() => onRemoveItem(item.id)}
                    style={{ background: '#e74c3c', color: '#fff', border: 'none', marginLeft: '4px' }}
                  >
                    ×
                  </button>
                </div>
              </li>
            ))}
          </ul>

          <div className="cart-total">
            Total: ${total.toFixed(2)}
          </div>

          <form className="customer-form" onSubmit={handleSubmit}>
            <h3>Datos de Entrega</h3>
            <div className="form-group">
              <label>Nombre:</label>
              <input
                type="text"
                value={customerName}
                onChange={(e) => setCustomerName(e.target.value)}
                placeholder="Ej: Juan Pérez"
                required
              />
            </div>
            <div className="form-group">
              <label>Teléfono (solo números):</label>
              <input
                type="tel"
                inputMode="numeric"
                pattern="[0-9]*"
                value={customerPhone}
                onChange={(e) => setCustomerPhone(e.target.value.replace(/\D/g, ''))}
                placeholder="Ej: 1122334455"
                required
              />
            </div>
            <div className="form-group">
              <label>Dirección:</label>
              <input
                type="text"
                value={customerAddress}
                onChange={(e) => setCustomerAddress(e.target.value)}
                placeholder="Ej: Av. Siempreviva 123"
                required
              />
            </div>

            <button
              type="submit"
              className="btn btn-success"
              style={{ marginTop: '0.5rem', width: '100%' }}
              disabled={loading || cartItems.length === 0}
            >
              {loading ? 'Procesando...' : 'Hacer Pedido por WhatsApp'}
            </button>
          </form>
        </>
      )}
    </div>
  );
}
