import React, { useState, useEffect } from 'react';
import { getPublicProducts } from '../api/products';
import { getCategories } from '../api/categories';
import { createOrder } from '../api/orders';
import ProductCard from '../components/ProductCard';
import Cart from '../components/Cart';

const WHATSAPP_NUMBER = import.meta.env.VITE_WHATSAPP_NUMBER || '5491112345678';

export default function Catalog() {
  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [cart, setCart] = useState(() => {
    try {
      const saved = localStorage.getItem('realico_cart');
      return saved ? JSON.parse(saved) : [];
    } catch (e) {
      return [];
    }
  });
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [successMessage, setSuccessMessage] = useState('');

  useEffect(() => {
    localStorage.setItem('realico_cart', JSON.stringify(cart));
  }, [cart]);

  const loadData = async () => {
    setLoading(true);
    setError('');
    try {
      const [prodsData, catsData] = await Promise.all([getPublicProducts(), getCategories()]);
      setProducts(prodsData || []);
      setCategories(catsData || []);
    } catch (err) {
      setError('No se pudo conectar con el backend. Asegúrate de que esté iniciado en el puerto 8080.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleAddToCart = (product) => {
    setCart((prevCart) => {
      const existing = prevCart.find((item) => item.id === product.id);
      if (existing) {
        if (existing.quantity >= product.stock) return prevCart;
        return prevCart.map((item) =>
          item.id === product.id ? { ...item, quantity: item.quantity + 1 } : item
        );
      }
      return [...prevCart, { ...product, quantity: 1 }];
    });
  };

  const handleUpdateQuantity = (productId, newQuantity) => {
    if (newQuantity <= 0) {
      handleRemoveItem(productId);
      return;
    }
    setCart((prevCart) =>
      prevCart.map((item) => (item.id === productId ? { ...item, quantity: newQuantity } : item))
    );
  };

  const handleRemoveItem = (productId) => {
    setCart((prevCart) => prevCart.filter((item) => item.id !== productId));
  };

  const handleSubmitOrder = async (orderPayload) => {
    setSubmitting(true);
    setError('');
    setSuccessMessage('');
    try {
      const createdOrder = await createOrder(orderPayload);

      // Armar mensaje de WhatsApp
      let msg = `*Nuevo Pedido #${createdOrder.id} - Realico Comidas*\n\n`;
      msg += `👤 *Cliente:* ${createdOrder.customer_name}\n`;
      msg += `📞 *Teléfono:* ${createdOrder.customer_phone}\n`;
      msg += `📍 *Dirección:* ${createdOrder.customer_address}\n\n`;
      msg += `📋 *Detalle del pedido:*\n`;

      cart.forEach((item) => {
        msg += `• ${item.quantity}x ${item.name} ($${(item.price * item.quantity).toFixed(2)})\n`;
      });

      msg += `\n💵 *Total:* $${createdOrder.total.toFixed(2)}`;

      const encodedMsg = encodeURIComponent(msg);
      const whatsappUrl = `https://wa.me/${WHATSAPP_NUMBER}?text=${encodedMsg}`;

      setCart([]);
      setSuccessMessage('¡Pedido confirmado correctamente! Redirigiendo a WhatsApp...');

      // Recargar catálogo actualizado (stocks descontados)
      loadData();

      // Abrir WhatsApp en nueva pestaña
      setTimeout(() => {
        window.open(whatsappUrl, '_blank');
      }, 1000);
    } catch (err) {
      setError(err.message || 'Error al procesar el pedido');
    } finally {
      setSubmitting(false);
    }
  };

  // Agrupar productos por categoría
  const productsByCategory = categories.map((category) => ({
    category,
    products: products.filter((p) => p.category_id === category.id),
  })).filter((group) => group.products.length > 0);

  // Productos sin categoría o categorías no listadas
  const uncategorizedProducts = products.filter(
    (p) => !categories.some((c) => c.id === p.category_id)
  );

  return (
    <div className="catalog-layout">
      <div className="catalog-main">
        <h2>Catálogo de Comidas</h2>
        <p style={{ marginBottom: '1.5rem', color: '#666' }}>
          Selecciona tus platos favoritos y realiza tu pedido directamente por WhatsApp.
        </p>

        {error && <div className="alert-error">{error}</div>}
        {successMessage && <div className="alert-success">{successMessage}</div>}

        {loading ? (
          <p>Cargando menú...</p>
        ) : products.length === 0 ? (
          <div className="table-container" style={{ textAlign: 'center', padding: '2rem' }}>
            <p>No hay productos disponibles en el menú en este momento.</p>
            <small style={{ color: '#777' }}>El administrador puede agregar categorías y productos desde el Panel Admin.</small>
          </div>
        ) : (
          <>
            {productsByCategory.map(({ category, products }) => (
              <div key={category.id} className="category-section">
                <h3 className="category-title">{category.name}</h3>
                <div className="products-grid">
                  {products.map((product) => {
                    const cartItem = cart.find((item) => item.id === product.id);
                    const cartQty = cartItem ? cartItem.quantity : 0;
                    return (
                      <ProductCard
                        key={product.id}
                        product={product}
                        cartQuantity={cartQty}
                        onAddToCart={handleAddToCart}
                      />
                    );
                  })}
                </div>
              </div>
            ))}

            {uncategorizedProducts.length > 0 && (
              <div className="category-section">
                <h3 className="category-title">Varios</h3>
                <div className="products-grid">
                  {uncategorizedProducts.map((product) => {
                    const cartItem = cart.find((item) => item.id === product.id);
                    const cartQty = cartItem ? cartItem.quantity : 0;
                    return (
                      <ProductCard
                        key={product.id}
                        product={product}
                        cartQuantity={cartQty}
                        onAddToCart={handleAddToCart}
                      />
                    );
                  })}
                </div>
              </div>
            )}
          </>
        )}
      </div>

      <Cart
        cartItems={cart}
        onUpdateQuantity={handleUpdateQuantity}
        onRemoveItem={handleRemoveItem}
        onSubmitOrder={handleSubmitOrder}
        loading={submitting}
      />
    </div>
  );
}
