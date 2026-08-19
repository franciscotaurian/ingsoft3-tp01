import React from 'react';

export default function Header({ activeTab, setActiveTab, isAdmin }) {
  return (
    <header>
      <h1>Realico Comidas</h1>
      <nav>
        <button
          className={activeTab === 'catalog' ? 'active' : ''}
          onClick={() => setActiveTab('catalog')}
        >
          Catálogo
        </button>

        {/* Historial de Pedidos solo visible para administradores autenticados */}
        {isAdmin && (
          <button
            className={activeTab === 'orders' ? 'active' : ''}
            onClick={() => setActiveTab('orders')}
          >
            Historial de Pedidos
          </button>
        )}

        <button
          className={activeTab === 'admin' ? 'active' : ''}
          onClick={() => setActiveTab('admin')}
        >
          Panel Admin
        </button>
      </nav>
    </header>
  );
}
