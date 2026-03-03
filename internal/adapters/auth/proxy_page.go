package auth

// successPage is the HTML served to the browser after cookies are captured.
const successPage = `<!DOCTYPE html>
<html lang="es">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>UA-CLI — Login exitoso</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      min-height: 100vh;
      display: flex;
      align-items: center;
      justify-content: center;
      background: linear-gradient(135deg, #0a1628 0%, #1a2744 50%, #0d3b66 100%);
      color: #e0e8f0;
    }
    .card {
      text-align: center;
      padding: 3rem 4rem;
      background: rgba(255,255,255,0.06);
      border-radius: 1.5rem;
      border: 1px solid rgba(255,255,255,0.1);
      backdrop-filter: blur(20px);
      box-shadow: 0 8px 32px rgba(0,0,0,0.3);
      max-width: 480px;
    }
    .check {
      font-size: 4rem;
      margin-bottom: 1rem;
      animation: pop 0.4s ease-out;
    }
    @keyframes pop {
      0% { transform: scale(0); opacity: 0; }
      80% { transform: scale(1.2); }
      100% { transform: scale(1); opacity: 1; }
    }
    h1 {
      font-size: 1.6rem;
      font-weight: 700;
      margin-bottom: 0.5rem;
      color: #4ade80;
    }
    p {
      font-size: 1rem;
      opacity: 0.7;
      line-height: 1.5;
    }
    .brand {
      margin-top: 2rem;
      font-size: 0.85rem;
      opacity: 0.4;
      letter-spacing: 0.1em;
    }
  </style>
</head>
<body>
  <div class="card">
    <div class="check">✓</div>
    <h1>Login exitoso</h1>
    <p>Tu sesión ha sido capturada.<br>Puedes cerrar esta pestaña.</p>
    <div class="brand">UA-CLI</div>
  </div>
</body>
</html>`
