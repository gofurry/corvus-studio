import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';

import App from './App';
import './styles/index.scss';

const rootElement = document.getElementById('root');

if (rootElement === null) {
  throw new Error('Corvus Studio root element is missing');
}

createRoot(rootElement).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
