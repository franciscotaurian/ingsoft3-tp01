import React, { useState, useEffect } from 'react';

export default function ProductForm({ categories, initialProduct, onSubmit, onCancel }) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [price, setPrice] = useState('');
  const [stock, setStock] = useState('');
  const [categoryId, setCategoryId] = useState('');
  const [imageUrl, setImageUrl] = useState('');

  useEffect(() => {
    if (initialProduct) {
      setName(initialProduct.name || '');
      setDescription(initialProduct.description || '');
      setPrice(initialProduct.price !== undefined ? initialProduct.price : '');
      setStock(initialProduct.stock !== undefined ? initialProduct.stock : '');
      setCategoryId(initialProduct.category_id || (categories[0]?.id || ''));
      setImageUrl(initialProduct.image_url || '');
    } else {
      setName('');
      setDescription('');
      setPrice('');
      setStock('');
      setCategoryId(categories[0]?.id || '');
      setImageUrl('');
    }
  }, [initialProduct, categories]);

  const handleSubmit = (e) => {
    e.preventDefault();
    onSubmit({
      name: name.trim(),
      description: description.trim(),
      price: parseFloat(price),
      stock: parseInt(stock, 10),
      category_id: parseInt(categoryId, 10),
      image_url: imageUrl.trim(),
    });
  };

  return (
    <form onSubmit={handleSubmit} style={{ background: '#f9f9f9', padding: '1rem', border: '1px solid #ddd', borderRadius: '6px', marginBottom: '1.5rem' }}>
      <h3>{initialProduct ? 'Editar Producto' : 'Nuevo Producto'}</h3>
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem', marginTop: '0.8rem' }}>
        <div className="form-group">
          <label>Nombre:</label>
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Ej: Hamburguesa Completa"
            required
          />
        </div>
        <div className="form-group">
          <label>Categoría:</label>
          <select value={categoryId} onChange={(e) => setCategoryId(e.target.value)} required>
            {categories.map((cat) => (
              <option key={cat.id} value={cat.id}>
                {cat.name}
              </option>
            ))}
          </select>
        </div>
        <div className="form-group">
          <label>Precio ($):</label>
          <input
            type="number"
            step="0.01"
            min="0.01"
            value={price}
            onChange={(e) => setPrice(e.target.value)}
            placeholder="Ej: 4500"
            required
          />
        </div>
        <div className="form-group">
          <label>Stock inicial:</label>
          <input
            type="number"
            min="0"
            value={stock}
            onChange={(e) => setStock(e.target.value)}
            placeholder="Ej: 20"
            required
          />
        </div>
        <div className="form-group" style={{ gridColumn: '1 / -1' }}>
          <label>Descripción:</label>
          <input
            type="text"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Ej: Carne 180g, queso cheddar, bacón, lechuga y tomate"
          />
        </div>
        <div className="form-group" style={{ gridColumn: '1 / -1' }}>
          <label>URL de Imagen (Opcional):</label>
          <input
            type="url"
            value={imageUrl}
            onChange={(e) => setImageUrl(e.target.value)}
            placeholder="https://ejemplo.com/imagen.jpg"
          />
        </div>
      </div>
      <div style={{ marginTop: '1rem', display: 'flex', gap: '0.5rem' }}>
        <button type="submit" className="btn btn-success">
          {initialProduct ? 'Guardar Cambios' : 'Crear Producto'}
        </button>
        {initialProduct && (
          <button type="button" className="btn btn-secondary" onClick={onCancel}>
            Cancelar
          </button>
        )}
      </div>
    </form>
  );
}
