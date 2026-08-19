import React, { useState, useEffect } from 'react';
import { getOrders } from '../api/orders';

export default function Orders() {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  // El historial de pedidos requiere credenciales de administrador
  const authHeader = localStorage.getItem('realico_admin_auth') || '';

  const loadOrders = async () => {
    setLoading(true);
    setError('');
    try {
      const data = await getOrders(authHeader);
      setOrders(data || []);
    } catch (err) {
      if (err.message && err.message.toLowerCase().includes('credencial')) {
        setError('Credenciales de administrador inválidas. Por favor, iniciá sesión desde el Panel Admin primero.');
      } else {
        setError(err.message || 'Error al cargar el historial de pedidos.');
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (!authHeader) {
      setError('El historial de pedidos es una sección privada. Debés iniciar sesión como administrador desde el Panel Admin primero.');
      setLoading(false);
      return;
    }
    loadOrders();
  }, [authHeader]);

  const formatDate = (dateStr) => {
    if (!dateStr) return '';
    const date = new Date(dateStr);
    return date.toLocaleString('es-AR');
  };

  return (
    <div>
      <h2>Historial de Pedidos</h2>
      <p style={{ marginBottom: '1.5rem', color: '#666' }}>
        Sección privada: solo accesible para administradores.
      </p>

      {error && (
        <div>
          <div className="alert-error">{error}</div>
          {!authHeader && (
            <p style={{ marginTop: '1rem', color: '#555' }}>
              Para acceder, dirigite al <strong>Panel Admin</strong> e iniciá sesión con tus credenciales.
            </p>
          )}
        </div>
      )}

      {loading && <p>Cargando historial...</p>}

      {!loading && !error && orders.length === 0 && (
        <div className="table-container" style={{ textAlign: 'center', padding: '2rem' }}>
          <p>No se registraron pedidos aún.</p>
        </div>
      )}

      {!loading && !error && orders.length > 0 && (
        <div className="table-container">
          <table>
            <thead>
              <tr>
                <th># Pedido</th>
                <th>Fecha</th>
                <th>Cliente</th>
                <th>Detalle</th>
                <th>Total</th>
                <th>Estado</th>
              </tr>
            </thead>
            <tbody>
              {orders.map((order) => (
                <tr key={order.id}>
                  <td>#{order.id}</td>
                  <td>{formatDate(order.created_at)}</td>
                  <td>
                    <strong>{order.customer_name}</strong>
                    <br />
                    <small>📞 {order.customer_phone}</small>
                    <br />
                    <small>📍 {order.customer_address}</small>
                  </td>
                  <td>
                    <ul style={{ paddingLeft: '1rem', margin: 0, fontSize: '0.9rem' }}>
                      {order.items?.map((item) => (
                        <li key={item.id}>
                          {item.quantity}x {item.product?.name || `Producto #${item.product_id}`} (${item.unit_price.toFixed(2)})
                        </li>
                      ))}
                    </ul>
                  </td>
                  <td><strong>${order.total.toFixed(2)}</strong></td>
                  <td>
                    <span className={`badge badge-${order.status}`}>
                      {order.status}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
