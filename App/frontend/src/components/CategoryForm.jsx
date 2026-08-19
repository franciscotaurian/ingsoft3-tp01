import React, { useState, useEffect } from 'react';

export default function CategoryForm({ initialCategory, onSubmit, onCancel }) {
  const [name, setName] = useState('');

  useEffect(() => {
    if (initialCategory) {
      setName(initialCategory.name);
    } else {
      setName('');
    }
  }, [initialCategory]);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!name.trim()) return;
    onSubmit(name.trim());
    setName('');
  };

  return (
    <form onSubmit={handleSubmit} style={{ display: 'flex', gap: '0.5rem', alignItems: 'center', marginBottom: '1rem' }}>
      <input
        type="text"
        value={name}
        onChange={(e) => setName(e.target.value)}
        placeholder="Nombre de categoría"
        style={{ padding: '0.5rem', borderRadius: '4px', border: '1px solid #ccc', flex: 1 }}
        required
      />
      <button type="submit" className="btn btn-success">
        {initialCategory ? 'Guardar' : 'Agregar Categoría'}
      </button>
      {initialCategory && (
        <button type="button" className="btn btn-secondary" onClick={onCancel}>
          Cancelar
        </button>
      )}
    </form>
  );
}
