import React from 'react';

export default function MetricsCards({ metrics }) {
  if (!metrics) return null;

  const pendingCount = metrics.orders_by_status?.pendiente || 0;
  const confirmedCount = metrics.orders_by_status?.confirmado || 0;
  const deliveredCount = metrics.orders_by_status?.entregado || 0;

  return (
    <div className="metrics-grid">
      <div className="metric-card">
        <h3>Total Facturado</h3>
        <div className="value">${(metrics.total_billed || 0).toFixed(2)}</div>
      </div>
      <div className="metric-card">
        <h3>Pedidos por Estado</h3>
        <div className="value" style={{ fontSize: '1.1rem' }}>
          ⏳ Pend: <strong>{pendingCount}</strong> | 🔵 Conf: <strong>{confirmedCount}</strong> | ✅ Entr: <strong>{deliveredCount}</strong>
        </div>
      </div>
      <div className="metric-card">
        <h3>Producto Más Vendido</h3>
        <div className="value" style={{ fontSize: '1.2rem' }}>
          {metrics.most_sold_product || 'N/A'}
        </div>
      </div>
      <div className="metric-card">
        <h3>Productos en Catálogo</h3>
        <div className="value">{metrics.total_products || 0}</div>
      </div>
      <div className="metric-card">
        <h3>Valor del Inventario</h3>
        <div className="value">${(metrics.inventory_value || 0).toFixed(2)}</div>
      </div>
    </div>
  );
}
