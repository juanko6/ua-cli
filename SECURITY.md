# 🔒 Política de Seguridad

## Versiones Soportadas

| Versión | Estado    |
|---------|-----------|
| Latest  | ✅ Soportada |
| Anterior | ❌ No soportada |

## Reportar una Vulnerabilidad

**Por favor NO abras un Issue público para reportar vulnerabilidades de seguridad.**

Si descubres una vulnerabilidad de seguridad en `ua-cli`, por favor sigue este proceso privado:

1. Abre un [Security Advisory privado](https://github.com/juanko6/ua-cli/security/advisories/new) en GitHub.
2. Incluye un descripción detallada del problema, los pasos para reproducirlo y el impacto potencial.
3. Recibirás una respuesta en un plazo máximo de **5 días laborables**.

## Contexto de Seguridad

`ua-cli` maneja **cookies de sesión de UACloud**. Tenga en cuenta que:

- Las cookies se almacenan localmente en el keyring del sistema operativo o en un archivo cifrado local (`~/.ua-cli/`).
- El CLI **nunca** envía tus credenciales a terceros.
- El CLI **nunca** loguea datos sensibles como cookies.
- El tráfico va siempre directamente a los servidores oficiales de la Universidad de Alicante (`*.ua.es`).

## Ámbito (Scope)

Las siguientes áreas son de especial interés para la seguridad:

- Almacenamiento y lectura de cookies de sesión.
- El proxy local de captura de cookies (`internal/adapters/auth/proxy.go`).
- Cualquier funcionalidad que transmita datos a servidores externos.
