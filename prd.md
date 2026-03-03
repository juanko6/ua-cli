# 📘 PRODUCT REQUIREMENTS DOCUMENT (PRD)
# ua-cli — University of Alicante Command-Line Companion

---

# 1. VISIÓN DEL PRODUCTO

**ua-cli** es la herramienta CLI definitiva para estudiantes de la Universidad de Alicante que transforma la experiencia lenta y web-dependiente de UACloud en interacciones instantáneas, composables y optimizadas para terminal.

Encapsula la filosofía Unix:  
> pequeños comandos, predecibles, rápidos, componibles.

---

# 2. PROBLEMA

El portal UACloud:
- Requiere múltiples clics y navegación lenta.
- No está optimizado para usuarios técnicos.
- No permite automatización.
- Tiene latencia perceptible.
- No ofrece una experiencia CLI ni scripting.

Los estudiantes técnicos pierden tiempo consultando:
- Horarios
- Notas
- Avisos
- Recursos del campus virtual

---

# 3. OBJETIVOS

## 🎯 Objetivo Principal (North Star Metric)

Reducir el tiempo de consulta de información académica a menos de 1 segundo por comando.

---

## Objetivos Secundarios

- Permitir consulta completa desde terminal.
- Minimizar logins repetidos.
- Garantizar almacenamiento seguro de credenciales.
- Mantener arquitectura mantenible y testeable.
- Evitar fricción de UX en entornos SSH.

---

# 4. ALCANCE (MVP)

## ✅ Incluido

1. Autenticación segura con persistencia de sesión.
2. Consulta de horario semanal (CLI).
3. Consulta de notas.
4. Consulta de avisos UACloud.
5. Descarga automática de recursos del campus virtual.
6. Soporte JSON output.
7. Soporte multi-OS.
8. Auto-update.
9. Disclaimer legal.
10. Soporte 2FA si la UA lo implementa.

---

## ❌ Excluido (por ahora)

- Modo offline completo.
- Push notifications del sistema.
- Organización avanzada de archivos.
- Multi-rol.
- Telemetría.
- Daemon en background.

---

# 5. USUARIOS

## Usuario Primario

Estudiantes de la Universidad de Alicante con acceso a UACloud.

## Futuro

Extensible a profesorado.

---

# 6. PRINCIPIOS FUNDAMENTALES

1. Speed is a Feature.
2. UX terminal-first.
3. Seguridad local-first.
4. Arquitectura hexagonal.
5. Scraping resiliente ante cambios estructurales.

---

# 7. REQUISITOS FUNCIONALES

## 7.1 Autenticación

- Login persistente.
- Uso de OS Keyring nativo.
- Soporte 2FA.
- Renovación automática de sesión cuando sea posible.
- Manejo robusto de expiración.

---

## 7.2 Schedule

- Mostrar horario semanal.
- Output en formato:
  - Tabla CLI
  - JSON opcional
- Alertar cambios de aula/horario respecto a último fetch.

---

## 7.3 Grades

- Mostrar listado de asignaturas y notas.
- Detectar nuevas notas publicadas.
- Output JSON.

---

## 7.4 Notices

- Listar avisos.
- Ordenados por fecha.
- Indicar nuevos desde último fetch.

---

## 7.5 Virtual Campus & Moodle

- **Material Docente**: Descarga automática y masiva de recursos.
- **Entregas (Deadlines)**: Integración con Moodle paralistar próximos trabajos, prácticas y fechas de cierre.
- Opciones de descarga directa y listado interactivo.

---

## 7.6 Tutorías e Interacción

- Enviar solicitudes de tutoría a un profesor rápidamente.
- Automatizar el flujo de: (Seleccionar asignatura) -> (Seleccionar profesor) -> (Enviar mensaje/solicitud) sin navegar la UI web de UACloud.

---

## 7.7 CLI Design

Ejemplos:

ua schedule
ua grades
ua notices
ua campus download
ua schedule --json

---

# 8. REQUISITOS NO FUNCIONALES

## Rendimiento
- < 1 segundo para comandos cacheados.
- Timeout configurable.

## Seguridad
- Credenciales solo en keyring.
- No logs con datos sensibles.
- No telemetría.

## Arquitectura
- Hexagonal (Ports & Adapters).
- Separación:
  - CLI
  - Application
  - Domain
  - Infrastructure
- Test coverage alta para parsers HTML/JSON.

## Compatibilidad
- macOS
- Linux
- Windows

---

# 9. RIESGOS Y MITIGACIONES

## Riesgo 1: Cambio en estructura UACloud
Mitigación:
- Parsers desacoplados.
- Tests snapshot.
- Abstracción de cliente HTTP.

## Riesgo 2: Bloqueo institucional
Mitigación:
- Rate limiting.
- User-Agent configurable.
- Disclaimer legal claro.

## Riesgo 3: Rotura por 2FA
Mitigación:
- Soporte interactivo OTP.
- Manejo de cookies seguro.

---

# 10. DISCLAIMER LEGAL

- Proyecto no oficial.
- No afiliado a la Universidad de Alicante.
- Uso bajo responsabilidad del usuario.
- Puede dejar de funcionar ante cambios estructurales.

---

# 11. DISTRIBUCIÓN

- Homebrew
- APT
- Scoop
- Binarios standalone
- Auto-update integrado

---

# 12. MÉTRICAS DE ÉXITO

| Métrica | Objetivo |
|----------|----------|
| Tiempo promedio de comando | <1s |
| Nº comandos promedio por sesión | 8 |
| Tasa de fallo por cambio estructural | <5% |

---

# 13. ROADMAP PROPUESTO

## Fase 1 – Core
- Login persistente
- Schedule
- Grades

## Fase 2 – Notificaciones
- Alertas cambios horario
- Alertas nuevas notas

## Fase 3 – Campus & Moodle
- Parseo de fechas de entregas (Assignments)
- Descarga masiva del material docente

## Fase 4 – Interacción & Tutorías
- Flujo interactivo para redactar correos/tutorías al profesorado
- Listado de profesores por asignatura

## Fase 5 – Bots & Notificaciones (Producción Full)
- Despliegue de bots (Telegram, Discord, WhatsApp)
- Alertas en tiempo real (push)
- Hardening & CI/CD avanzado

---

# 14. CRITERIOS DE ACEPTACIÓN GLOBALES

- Todos los comandos funcionan sin TUI.
- Compatible con piping Unix.
- Soporte JSON en todas las consultas.
- Manejo elegante de errores.
- No requiere login frecuente.
- Instalación < 1 minuto.

---

# 15. FUTURAS EXTENSIONES Y BOTS INTERACTIVOS

El proyecto evolucionará de una mera CLI a un motor de automatización que exponga estas facilidades mediante bots para uso inmediato on-the-go:

- **Plataformas**: Telegram, WhatsApp, Discord, y Slack.
- **Funcionamiento**: Añadir el bot a un chat, autenticar con un solo comando, y solicitar:
  - `/horario`
  - `/entregas_pendientes`
  - `/avisos`
- **Tutorías automatizadas**: Notificaciones push integradas cuando el profesor responda.
- **Profesores**: Soporte futuro para perfil de usuario profesor.

---

# 16. DEFINICIÓN DE DONE

- Binarios multiplataforma.
- Tests automatizados.
- Documentación CLI.
- Disclaimer visible.
- Versionado semántico.
- Changelog.

---

# RESUMEN EJECUTIVO

ua-cli será una herramienta OSS personal que:
- Optimiza tiempo
- Aumenta productividad técnica
- Respeta privacidad
- Mantiene arquitectura limpia
- Minimiza fricción universitaria

