import React from 'react';

export default function ProductCard({ product, onAddToCart, cartQuantity }) {
  const isOutOfStock = product.stock <= 0 || cartQuantity >= product.stock;

  return (
    <div className="product-card">
      <div>
        {product.image_url ? (
          <img src={product.image_url} alt={product.name} className="product-img" />
        ) : (
          <div className="product-img" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#888' }}>
            Sin Imagen
          </div>
        )}
        <div className="product-title">{product.name}</div>
        <div className="product-description">{product.description || 'Sin descripción'}</div>
      </div>
      <div>
        <div className="product-price">${product.price.toFixed(2)}</div>
        <div className="product-stock">Stock: {product.stock - cartQuantity} disp.</div>
        <button
          className="btn btn-success"
          style={{ width: '100%' }}
          onClick={() => onAddToCart(product)}
          disabled={isOutOfStock}
        >
          {isOutOfStock ? 'Sin Stock' : 'Agregar al Carrito'}
        </button>
      </div>
    </div>
  );
}
