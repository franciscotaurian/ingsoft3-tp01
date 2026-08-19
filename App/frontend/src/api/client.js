const API_URL = import.meta.env.VITE_API_URL || '';

export async function fetchApi(endpoint, options = {}) {
  const url = `${API_URL}${endpoint}`;
  const headers = {
    'Content-Type': 'application/json',
    ...(options.headers || {}),
  };

  const config = {
    ...options,
    headers,
  };

  const response = await fetch(url, config);

  if (!response.ok) {
    let errorMessage = 'Error en la petición';
    try {
      const data = await response.json();
      if (data.error) {
        errorMessage = data.error;
      }
    } catch (e) {
      // Ignorar parse error
    }
    throw new Error(errorMessage);
  }

  // Para respuestas 204 No Content o vacías
  if (response.status === 204) {
    return null;
  }

  return response.json();
}
