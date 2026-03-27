# 🤝 Contributing to ua-cli

¡Gracias por tu interés en contribuir a **ua-cli**! Este proyecto está construido por y para estudiantes de la Universidad de Alicante. Toda ayuda es bienvenida.

---

## 📋 Tabla de Contenidos

- [Código de Conducta](#código-de-conducta)
- [¿Cómo Puedo Contribuir?](#cómo-puedo-contribuir)
- [Metodología de Desarrollo (SDD)](#metodología-de-desarrollo-sdd)
- [Configuración del Entorno](#configuración-del-entorno)
- [Estilo de Código y Convenciones](#estilo-de-código-y-convenciones)
- [Flujo de Pull Requests](#flujo-de-pull-requests)
- [Commit Convention](#commit-convention)

---

## Código de Conducta

Este proyecto aplica un entorno respetuoso e inclusivo para todos. Al participar, aceptas cumplir con el [Código de Conducta](CODE_OF_CONDUCT.md).

---

## ¿Cómo Puedo Contribuir?

### 🐛 Reportar un Bug

Si encuentras un error, abre un [Issue](https://github.com/juanko6/ua-cli/issues/new?template=bug_report.yml) usando la plantilla **Bug Report**. Por favor, incluye:
- Versión del CLI (`ua --version`)
- Sistema operativo y versión
- Pasos exactos para reproducirlo
- Salida completa del error (copia el texto o un screenshot)

### 💡 Proponer una Nueva Feature

1. Asegúrate de que no existe ya un Issue o PR similar abierto.
2. Abre un [Feature Request](https://github.com/juanko6/ua-cli/issues/new?template=feature_request.yml) describiendo el caso de uso.
3. Espera feedback del equipo antes de empezar a codificar (podemos rechazar features que no encajen con la visión del [PRD](prd.md)).

### 🛠️ Escribir Código

¡Perfecto! Sigue la [Metodología de Desarrollo (SDD)](#metodología-de-desarrollo-sdd) para que tu contribución sea aprobada lo más rápido posible.

---

## Metodología de Desarrollo (SDD)

**ua-cli** usa **Specification-Driven Development (SDD)** mediante los workflows de [Speckit](.agent/workflows/). Toda nueva feature **debe** pasar por el siguiente ciclo antes de escribir código:

```
speckit.specify → speckit.plan → speckit.tasks → speckit.implement
```

| Paso | Workflow | Output |
|------|----------|--------|
| 1. Especificación | `/speckit.specify` | `specs/<feature>/spec.md` |
| 2. Diseño Técnico | `/speckit.plan` | `plan.md`, `research.md`, `data-model.md` |
| 3. Tareas Atómicas | `/speckit.tasks` | `tasks.md` |
| 4. Implementación | `/speckit.implement` | Código, tests, PR |

Mira la carpeta `specs/` para ver ejemplos reales del proceso.

> ⚠️ Los PRs que añadan código sin documentación de diseño en `specs/` pueden ser rechazados por falta de contexto.

---

## Configuración del Entorno

### Prerequisitos

- **Go 1.24+** — https://go.dev/dl/
- **Git**

### Setup

```bash
# 1. Fork y clona el repositorio
git clone https://github.com/<TU_USUARIO>/ua-cli.git
cd ua-cli

# 2. Descarga dependencias
go mod download

# 3. Compila y verifica
go build -o ua ./cmd/ua-cli/
./ua --help

# 4. Ejecuta los tests
go test ./...
```

### Configurar el Linter

Este proyecto usa `golangci-lint`. Puedes [instalarlo aquí](https://golangci-lint.run/welcome/install/).

```bash
golangci-lint run
```

---

## Estilo de Código y Convenciones

- **Formato**: Todo el código debe pasar `go fmt ./...` sin diferencias.
- **Arquitectura**: Seguimos los principios de **Ports & Adapters** (Arquitectura Hexagonal). El código de dominio NO debe importar adaptadores.
- **Errores**: Siempre chequea y maneja los errores retornados. No uses `_ = err` a menos que sea intencional y documentado.
- **Tests**: Escribe tests unitarios para los parsers, servicios y lógica de dominio.
- **Idioma**: Los comentarios de código y variables pueden estar en inglés o español. La documentación de usuario en español.

### Arquitectura

```
cmd/ua-cli/       → Comandos Cobra (puntos de entrada)
internal/
  domain/         → Entidades y puertos (interfaces)
  service/        → Casos de uso
  adapters/       → Implementaciones concretas
specs/            → Documentación SDD por feature
```

---

## Flujo de Pull Requests

1. Crea una rama desde `main` con el formato `<numero>-<descripcion-corta>`:
   ```bash
   git checkout -b 003-ua-grades-command
   ```

2. Sigue el ciclo SDD (si aplica) y desarrolla tu feature.

3. Asegúrate de que todo el pipeline local pase antes de abrir el PR:
   ```bash
   go fmt ./...
   go test ./...
   golangci-lint run
   ```

4. Abre el PR con el título en formato [`Conventional Commits`](https://www.conventionalcommits.org/):
   - `feat: add ua grades command`
   - `fix: handle expired session cookie gracefully`
   - `docs: update contributing guidelines`

5. Completa la plantilla del PR y marca todos los checks.

6. Espera revisión. Todos los PRs requieren al menos **1 aprobación** de un maintainer.

---

## Commit Convention

Usamos [Conventional Commits](https://www.conventionalcommits.org/):

| Prefijo | Cuándo usarlo |
|---------|--------------|
| `feat:` | Nueva funcionalidad |
| `fix:` | Corrección de bug |
| `docs:` | Cambios en documentación |
| `refactor:` | Refactors sin cambio funcional |
| `test:` | Añadir o modificar tests |
| `ci:` | Cambios en GitHub Actions / workflows |
| `build:` | Cambios en dependencias (`go.mod`, `go.sum`) |
| `chore:` | Tareas de mantenimiento menores |

---

*Desarrollado con ☕ para sobrevivir a la lentitud de los portales universitarios.*
