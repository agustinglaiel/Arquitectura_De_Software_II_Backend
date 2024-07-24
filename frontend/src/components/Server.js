// server.js
const express = require('express');
const cors = require('cors');
const app = express();

// Middleware CORS
app.use(cors());

// Middleware para parsear JSON
app.use(express.json());

// Ejemplo de una ruta
app.get('/api/some-endpoint', (req, res) => {
  res.json({ message: 'Hello from the server!' });
});

// Puerto en el que escucha tu aplicación
const PORT = process.env.PORT || 5000;
app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});
