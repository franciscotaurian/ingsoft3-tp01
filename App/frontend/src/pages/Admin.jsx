import React, { useState, useEffect } from 'react';
import { getAdminProducts, createProduct, updateProduct, deleteProduct } from '../api/products';
import { getCategories, createCategory, updateCategory, deleteCategory } from '../api/categories';
import { getOrders, updateOrderStatus, getAdminMetrics } from '../api/orders';

import MetricsCards from '../components/MetricsCards';
import CategoryForm from '../components/CategoryForm';
import ProductForm from '../components/ProductForm';

export default function Admin({ onAuthChange }) {
  const [user, setUser] = useState('');
  const [password, setPassword] = useState('');
  const [authHeader, setAuthHeader] = useState(() => localStorage.getItem('realico_admin_auth') || '');
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  const [metrics, setMetrics] = useState(null);
  const [categories, setCategories] = useState([]);
  const [products, setProducts] = useState([]);
  const [orders, setOrders] = useState([]);

  const [editingCategory, setEditingCategory] = useState(null);
  const [editingProduct, setEditingProduct] = useState(null);
  const [showProductForm, setShowProductForm] = useState(false);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [successMessage, setSuccessMessage] = useState('');

  const handleLogin = (e) => {
    e.preventDefault();
    const token = 'Basic ' + btoa(`${user}:${password}`);
    setAuthHeader(token);
    localStorage.setItem('realico_admin_auth', token);
  };

  const handleLogout = () => {
    setAuthHeader('');
    setIsAuthenticated(false);
    localStorage.removeItem('realico_admin_auth');
    if (onAuthChange) onAuthChange(false); // ocultar Historial en navbar
  };

  const loadAdminData = async () => {
    if (!authHeader) return;
    setLoading(true);
    setError('');
    try {
      const [metData, catsData, prodsData, ordsData] = await Promise.all([
        getAdminMetrics(authHeader),
        getCategories(),
        getAdminProducts(authHeader),
        getOrders(authHeader),
      ]);


      setMetrics(metData);
      setCategories(catsData || []);
      setProducts(prodsData || []);
      setOrders(ordsData || []);
      setIsAuthenticated(true);
      if (onAuthChange) onAuthChange(true); // mostrar Historial en navbar
    } catch (err) {
      setError(err.message || 'Credenciales de administrador inválidas');
      setIsAuthenticated(false);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (authHeader) {
      loadAdminData();
    }
  }, [authHeader]);

  // --- Handlers de Categorías ---
  const handleCategorySubmit = async (name) => {
    setError('');
    setSuccessMessage('');
    try {
      if (editingCategory) {
        await updateCategory(editingCategory.id, name, authHeader);
        setSuccessMessage(`Categoría "${name}" actualizada.`);
        setEditingCategory(null);
      } else {
        await createCategory(name, authHeader);
        setSuccessMessage(`Categoría "${name}" creada.`);
      }
      loadAdminData();
    } catch (err) {
      setError(err.message);
    }
  };

  const handleCategoryDelete = async (id, name) => {
    if (!window.confirm(`¿Seguro que deseas eliminar la categoría "${name}"?`)) return;
    setError('');
    setSuccessMessage('');
    try {
      await deleteCategory(id, authHeader);
      setSuccessMessage(`Categoría "${name}" eliminada.`);
      loadAdminData();
    } catch (err) {
      setError(err.message);
    }
  };

  // --- Handlers de Productos ---
  const handleProductSubmit = async (productData) => {
    setError('');
    setSuccessMessage('');
    try {
      if (editingProduct) {
        await updateProduct(editingProduct.id, productData, authHeader);
        setSuccessMessage(`Producto "${productData.name}" actualizado.`);
        setEditingProduct(null);
      } else {
        await createProduct(productData, authHeader);
        setSuccessMessage(`Producto "${productData.name}" creado.`);
      }
      setShowProductForm(false);
      loadAdminData();
    } catch (err) {
      setError(err.message);
    }
  };

  const handleProductDelete = async (id, name) => {
    if (!window.confirm(`¿Seguro que deseas eliminar el producto "${name}"?`)) return;
    setError('');
    setSuccessMessage('');
    try {
      await deleteProduct(id, authHeader);
      setSuccessMessage(`Producto "${name}" eliminado.`);
      loadAdminData();
    } catch (err) {
      setError(err.message);
    }
  };

  // --- Handlers de Estado de Pedidos ---
  const handleStatusAdvance = async (orderId, currentStatus) => {
    let nextStatus = '';
    if (currentStatus === 'pendiente') nextStatus = 'confirmado';
    else if (currentStatus === 'confirmado') nextStatus = 'entregado';
    else return;

    setError('');
    setSuccessMessage('');
    try {
      await updateOrderStatus(orderId, nextStatus, authHeader);
      setSuccessMessage(`Pedido #${orderId} actualizado a "${nextStatus}".`);
      loadAdminData();
    } catch (err) {
      setError(err.message);
    }
  };

  if (!isAuthenticated) {
    return (
      <div className="login-box">
        <h2>Acceso de Administrador</h2>

        {error && <div className="alert-error">{error}</div>}

        <form onSubmit={handleLogin} className="customer-form">
          <div className="form-group">
            <label>Usuario:</label>
            <input
              type="text"
              value={user}
              onChange={(e) => setUser(e.target.value)}
              placeholder="admin"
              required
            />
          </div>
          <div className="form-group">
            <label>Contraseña:</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              required
            />
          </div>
          <button type="submit" className="btn btn-success" disabled={loading}>
            {loading ? 'Verificando...' : 'Ingresar'}
          </button>
        </form>
      </div>
    );
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
        <h2>Panel de Administración</h2>
        <button onClick={handleLogout} className="btn btn-secondary">
          Cerrar Sesión
        </button>
      </div>

      {error && <div className="alert-error">{error}</div>}
      {successMessage && <div className="alert-success">{successMessage}</div>}

      {/* 1. Dashboard de Métricas */}
      <MetricsCards metrics={metrics} />

      {/* 2. CRUD Categorías */}
      <div className="table-container">
        <h3>Gestión de Categorías</h3>
        <p style={{ fontSize: '0.85rem', color: '#666', marginBottom: '1rem' }}>
          Regla: No se puede eliminar una categoría que tenga productos asociados.
        </p>
        <CategoryForm
          initialCategory={editingCategory}
          onSubmit={handleCategorySubmit}
          onCancel={() => setEditingCategory(null)}
        />
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Nombre</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            {categories.length === 0 ? (
              <tr>
                <td colSpan="3" style={{ textAlign: 'center' }}>No hay categorías. Crea una arriba.</td>
              </tr>
            ) : (
              categories.map((cat) => (
                <tr key={cat.id}>
                  <td>{cat.id}</td>
                  <td><strong>{cat.name}</strong></td>
                  <td>
                    <button
                      className="btn"
                      style={{ padding: '0.3rem 0.6rem', fontSize: '0.8rem', marginRight: '0.4rem' }}
                      onClick={() => setEditingCategory(cat)}
                    >
                      Editar
                    </button>
                    <button
                      className="btn btn-danger"
                      style={{ padding: '0.3rem 0.6rem', fontSize: '0.8rem' }}
                      onClick={() => handleCategoryDelete(cat.id, cat.name)}
                    >
                      Eliminar
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* 3. CRUD Productos */}
      <div className="table-container">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
          <h3>Gestión de Productos</h3>
          {!showProductForm && !editingProduct && (
            <button
              className="btn btn-success"
              onClick={() => {
                if (categories.length === 0) {
                  alert('Primero debes crear al menos una categoría.');
                  return;
                }
                setShowProductForm(true);
              }}
            >
              + Nuevo Producto
            </button>
          )}
        </div>

        {(showProductForm || editingProduct) && (
          <ProductForm
            categories={categories}
            initialProduct={editingProduct}
            onSubmit={handleProductSubmit}
            onCancel={() => {
              setShowProductForm(false);
              setEditingProduct(null);
            }}
          />
        )}

        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Producto</th>
              <th>Categoría</th>
              <th>Precio</th>
              <th>Stock</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            {products.length === 0 ? (
              <tr>
                <td colSpan="6" style={{ textAlign: 'center' }}>No hay productos registrados.</td>
              </tr>
            ) : (
              products.map((prod) => (
                <tr key={prod.id}>
                  <td>{prod.id}</td>
                  <td>
                    <strong>{prod.name}</strong>
                    <br />
                    <small style={{ color: '#666' }}>{prod.description}</small>
                  </td>
                  <td>{prod.category?.name || `ID ${prod.category_id}`}</td>
                  <td>${prod.price.toFixed(2)}</td>
                  <td>
                    <span style={{ color: prod.stock === 0 ? 'red' : 'inherit', fontWeight: prod.stock === 0 ? 'bold' : 'normal' }}>
                      {prod.stock}
                    </span>
                  </td>
                  <td>
                    <button
                      className="btn"
                      style={{ padding: '0.3rem 0.6rem', fontSize: '0.8rem', marginRight: '0.4rem' }}
                      onClick={() => {
                        setEditingProduct(prod);
                        setShowProductForm(false);
                      }}
                    >
                      Editar
                    </button>
                    <button
                      className="btn btn-danger"
                      style={{ padding: '0.3rem 0.6rem', fontSize: '0.8rem' }}
                      onClick={() => handleProductDelete(prod.id, prod.name)}
                    >
                      Eliminar
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* 4. Lista de Pedidos & Transición de Estado */}
      <div className="table-container">
        <h3>Control de Pedidos y Transición de Estados</h3>
        <p style={{ fontSize: '0.85rem', color: '#666', marginBottom: '1rem' }}>
          Flujo permitido: Pendiente → Confirmado → Entregado
        </p>
        <table>
          <thead>
            <tr>
              <th># ID</th>
              <th>Cliente / Contacto</th>
              <th>Detalle de Productos</th>
              <th>Total</th>
              <th>Estado Actual</th>
              <th>Acción de Estado</th>
            </tr>
          </thead>
          <tbody>
            {orders.length === 0 ? (
              <tr>
                <td colSpan="6" style={{ textAlign: 'center' }}>No hay pedidos recibidos.</td>
              </tr>
            ) : (
              orders.map((ord) => (
                <tr key={ord.id}>
                  <td>#{ord.id}</td>
                  <td>
                    <strong>{ord.customer_name}</strong>
                    <br />
                    <small>📞 {ord.customer_phone}</small>
                    <br />
                    <small>📍 {ord.customer_address}</small>
                  </td>
                  <td>
                    <ul style={{ paddingLeft: '1rem', margin: 0, fontSize: '0.85rem' }}>
                      {ord.items?.map((item) => (
                        <li key={item.id}>
                          {item.quantity}x {item.product?.name || `Producto #${item.product_id}`} (${item.unit_price.toFixed(2)})
                        </li>
                      ))}
                    </ul>
                  </td>
                  <td><strong>${ord.total.toFixed(2)}</strong></td>
                  <td>
                    <span className={`badge badge-${ord.status}`}>
                      {ord.status}
                    </span>
                  </td>
                  <td>
                    {ord.status === 'pendiente' && (
                      <button
                        className="btn"
                        style={{ padding: '0.4rem 0.8rem', fontSize: '0.85rem' }}
                        onClick={() => handleStatusAdvance(ord.id, ord.status)}
                      >
                        Confirmar Pedido ➔
                      </button>
                    )}
                    {ord.status === 'confirmado' && (
                      <button
                        className="btn btn-success"
                        style={{ padding: '0.4rem 0.8rem', fontSize: '0.85rem' }}
                        onClick={() => handleStatusAdvance(ord.id, ord.status)}
                      >
                        Marcar Entregado ✓
                      </button>
                    )}
                    {ord.status === 'entregado' && (
                      <span style={{ color: '#27ae60', fontSize: '0.85rem', fontWeight: 'bold' }}>✓ Finalizado</span>
                    )}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
