## 📋 Descripción

<!--
Describe brevemente QUÉ cambia este PR y POR QUÉ es necesario.
- ¿Qué problema resuelve?
- ¿Qué nueva funcionalidad añade?
- Incluye el número de issue relacionado si existe: "Closes #123"
-->

Closes #

---

## 🔗 Tipo de Cambio

<!-- Marca con [x] los que apliquen: -->

- [ ] 🐛 `fix` - Corrección de bug
- [ ] ✨ `feat` - Nueva funcionalidad
- [ ] 📝 `docs` - Actualización de documentación
- [ ] ♻️ `refactor` - Refactorización sin cambio funcional
- [ ] 🧪 `test` - Añadir o modificar tests
- [ ] ⚙️ `ci` - Cambios en CI/CD
- [ ] 🔧 `build` - Cambios en dependencias

---

## 📐 Diseño SDD (Specification-Driven Development)

<!-- Si este PR introduce una nueva feature, enlaza los artefactos creados del ciclo SDD: -->

- [ ] No aplica (bug fix, docs, refactor)
- [ ] **spec.md**: `specs/<feature>/spec.md`
- [ ] **plan.md**: `specs/<feature>/plan.md`
- [ ] **tasks.md**: `specs/<feature>/tasks.md`

---

## ✅ Checklist

Antes de enviar el PR, verifica que:

- [ ] `go fmt ./...` no produce diferencias.
- [ ] `go test ./...` pasa sin errores.
- [ ] `golangci-lint run` no reporta issues.
- [ ] `go build -o ua ./cmd/ua-cli/` compila correctamente.
- [ ] Los cambios son compatibles con el estilo de arquitectura hexagonal del proyecto.
- [ ] Si aplica, los comandos nuevos tienen soporte `--json`.
- [ ] Leí el [CONTRIBUTING.md](../CONTRIBUTING.md).

---

## 📸 Capturas / Demo

<!-- Si hay cambios en la UI del terminal, añade una captura aquí: -->

