<div align="center">
  <h1>🎓 ua-cli</h1>
  <p><strong>La herramienta CLI inteligente para estudiantes de la Universidad de Alicante.</strong></p>
  <p>Transforma la experiencia lenta y web-dependiente de UACloud en comandos instantáneos, fluidos y optimizados para terminal.</p>

  [![Go Reference](https://pkg.go.dev/badge/github.com/juanko6/ua-cli.svg)](https://pkg.go.dev/github.com/juanko6/ua-cli)
  [![Latest Release](https://img.shields.io/github/v/release/juanko6/ua-cli?color=blue)](https://github.com/juanko6/ua-cli/releases)
  [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

---

## 🚀 Características (v0.3.x)

*   **🔑 Smart Auto-Login (CAS SSO):** `ua login` levanta un proxy local, abre tu navegador al portal de la universidad y captura transparente y automáticamente las cookies de sesión (sin copiar y pegar).
*   **📅 Horario Interactivo (TUI):** `ua schedule` recupera y parsea tu horario semanal real, mostrando asignaturas, horarios y **localidad de aula** (p.e. `A3/0007`) usando una interfaz interactiva de terminal o formato de tabla Unix.
*   **⚡ Ultra-rápido:** Tiempo de consulta sub-segundo garantizado.
*   **🛡️ Arquitectura Segura:** Las sesiones (cookies) se almacenan localmente y encriptadas por el sistema operativo. Arquitectura hexagonal clara con tests rigurosos.
*   **JSON API Ready:** Salida compatible con JSON (`--json`) para construir pipelines y automatizar flujos.
*   **Cross-Platform:** Soporte nativo y probado en Windows, macOS y Linux.

## 🛠️ Stack Tecnológico

*   **Lenguaje:** [Go (Golang)](https://go.dev/) 1.25+ - Para compilación a binarios nativos estáticos.
*   **CLI Framework:** [Cobra](https://github.com/spf13/cobra) - Para crear un CLI potente y extensible con subcomandos.
*   **Renderizado UI / TUI:** [Bubbletea](https://github.com/charmbracelet/bubbletea) & [Lipgloss](https://github.com/charmbracelet/lipgloss) & [Bubbles](https://github.com/charmbracelet/bubbles) - Para UI en la terminal interactiva de clase mundial.
*   **Proxy & Networking:** Interfaces HTTP Nativas de Go (`net/http/httputil`).
*   **CI/CD:** [GoReleaser](https://goreleaser.com/) & GitHub Actions - Para empaquetado y releases automatizados de binarios.

## 📦 Instalación

### Descarga del Binario (Multi-plataforma)
Ve a la página de [**Releases**](https://github.com/juanko6/ua-cli/releases) y descarga el ejecutable correspondiente a tu sistema operativo (Windows, macOS o Linux).

Extrae el archivo y coloca el binario (`ua` o `ua.exe`) en una carpeta que esté en tu `PATH` del sistema.

### Desde Código Fuente (Requiere Go)
```bash
git clone https://github.com/juanko6/ua-cli.git
cd ua-cli
go build -o ua ./cmd/ua-cli/
```

## 💻 Guía de Uso Rápido

### 1. Iniciar sesión en UACloud

```bash
ua login
```
El CLI abrirá tu navegador por defecto. Inicia sesión en la página oficial (Auténtica / CAS) y vuelve a tu terminal. El CLI guardará tu sesión de forma automática y segura. ¡Listo para la acción!

> *¿Estás en un servidor SSH sin navegador?* Puedes usar el flag de retorno: `ua login --cookie` y pegar el string manualmente.

Puedes revisar la salud de tu sesión en cualquier momento:
```bash
ua login --status
```

### 2. Consultar Horario 

```bash
# Ver horario de la semana actual
ua schedule

# Navegar por semanas
ua schedule --next
ua schedule --prev

# Output en JSON crudo para automatización (JQ, etc)
ua schedule --json
```

## 🏗️ Estructura del Proyecto (Arquitectura Hexagonal)

El proyecto está diseñado bajo los principios arquitectónicos de *Ports & Adapters*:

```text
├── cmd/ua-cli/               # Punto de entrada de la aplicación CLI (Cobra config)
├── internal/
│   ├── domain/               # Entidades de negocio (Schedule, Auth)
│   ├── service/              # Casos de uso de la aplicación principal
│   └── adapters/             # Implementaciones concretas de puertos:
│       ├── auth/             #   Proxy local, CAS login form handler y storage
│       ├── presenter/        #   Vistas TUI/Bubbletea y Text tables
│       ├── repo/             #   Almacenamiento caché local
│       └── uacloud/          #   Scrapers y API JSON para recuperar datos UA cloud
├── specs/                    # Documentación de diseño, SpecKit scripts y workflows
└── .github/workflows/        # Pipelines CI/CD y automatización GoReleaser
```

## 🤝 Roadmap / Futuro

Consulta el [Product Requirements Document (PRD)](prd.md) completo de este repositorio. Las funcionalidades siguientes en la mira son:

*   **Fase 1 (Core):** Login persistente (`ua login`), Horario semanal (`ua schedule`), y la clase actual (`ua now`).
*   **Fase 2 (Notificaciones):** Comando general para avisos del campus y cambios de aula (`ua notices`). Consulta de notas interactivas (`ua grades`).
*   **Fase 3 (Campus & Moodle):** Integración de Moodle para fechas de entregas (Deadlines) y descarga masiva de Material Docente.
*   **Fase 4 (Interacción & Tutorías):** Flujo y wizard integrado en terminal para enviar consultas de tutorías a profesores automatizadamente.
*   **Fase 5 (Bots & Notificaciones):** Despliegue de la API y comandos como un Chatbot nativo cross-platform (Telegram, Discord, WhatsApp).

## 📐 Specification-Driven Development (SDD) con Speckit

El ciclo de desarrollo de **`ua-cli`** no se hace escribiendo código directamente, sino que seguimos un modelo riguroso de diseño llamado *"Spec-Driven Development"* (SDD). Usamos los flujos de trabajo locales de la carpeta `.agent/workflows/` (Speckit) para planificar, validar e implementar funciones.

Si vas a proponer o crear una funcionalidad nueva, el ciclo de vida estricto es el siguiente:

1. **`@/speckit.specify`**: A partir de lenguaje natural y el PRD, se genera el documento formal de especificación `spec.md` detallando historias de usuario, requerimientos funcionales y una checklist de éxito en la carpeta de la nueva feature.
2. **`@/speckit.plan`**: Creamos los documentos técnicos (`plan.md`, `research.md`, y de modelo de datos `data-model.md`). El sistema propone cómo integrarlo en la arquitectura hexagonal del código de Go sin romper nada y con validación técnica.
3. **`@/speckit.tasks`**: Partiendo del conocimiento y el diseño anterior, se genera de forma exacta la subdivisión del listado de tareas atómicas interconectadas en el documento `tasks.md`.
4. **`@/speckit.implement`**: Arrancamos el motor para codificar tarea por tarea utilizando un enfoque estructurado (TDD si aplica), pasando por los pasos de infraestructura, funcionalidad base y acabados estéticos con `lipgloss` hasta el último paso del Pipeline CI/CD.

Siéntete libre de navegar a la carpeta `specs/` para ver los ejemplos exactos del diseño previo al código.

## ⚠️ Disclaimer Legal

**ua-cli** es un proyecto *Open Source* **NO OFICIAL** de terceros, desarrollado de forma independiente por estudiantes técnicos.
No está afiliado, endorsado ni soportado por la Universidad de Alicante.
Utiliza este software bajo tu propia responsabilidad. Funciona conectándose a los puntos de red públicos expuestos para los estudiantes en los portales UA. Su uso y la gestión local que este realiza de las cookies cumple con proteger a los estudiantes.

---

*Desarrollado con ☕ para sobrevivir a la lentitud de los portales universitarios tradicionales.*
