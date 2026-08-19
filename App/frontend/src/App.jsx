import React, { useState } from 'react';
import Header from './components/Header';
import Catalog from './pages/Catalog';
import Orders from './pages/Orders';
import Admin from './pages/Admin';

export default function App() {
  const [activeTab, setActiveTab] = useState('catalog');
  // Estado de autenticación compartido: lo lee el Header para mostrar/ocultar
  // el botón de Historial, y Admin lo actualiza al hacer login/logout.
  const [isAdmin, setIsAdmin] = useState(
    () => !!localStorage.getItem('realico_admin_auth')
  );

  return (
    <div>
      <Header
        activeTab={activeTab}
        setActiveTab={setActiveTab}
        isAdmin={isAdmin}
      />
      <main className="container">
        {activeTab === 'catalog' && <Catalog />}
        {activeTab === 'orders' && <Orders />}
        {activeTab === 'admin' && (
          <Admin onAuthChange={setIsAdmin} />
        )}
      </main>
    </div>
  );
}
