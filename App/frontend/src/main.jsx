import React from 'react';
import ReactDOM from 'react-dom/client';
import * as Sentry from '@sentry/react';
import App from './App.jsx';
import './styles.css';

Sentry.init({
  dsn: 'https://4d47de6a57b822afa194867970facf03@o4512286161043456.ingest.us.sentry.io/4512286185684992',
  environment: window.location.hostname, // 'localhost', 'miapp-front-prod...'
  integrations: [Sentry.browserTracingIntegration()],
  tracesSampleRate: 1.0, // medir 100% de las operaciones para la demo
});

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);