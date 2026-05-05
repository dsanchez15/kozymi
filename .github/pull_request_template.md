## Descripción

Breve descripción de qué hace este PR y por qué es necesario.

Fixes # (número de issue si aplica)

## Tipo de cambio

- [ ] `feat`: Nueva funcionalidad
- [ ] `fix`: Corrección de bug
- [ ] `docs`: Cambios en documentación
- [ ] `style`: Formato (sin cambios de lógica)
- [ ] `refactor`: Refactorización de código
- [ ] `perf`: Mejora de performance
- [ ] `test`: Tests
- [ ] `chore`: Tareas de mantenimiento

## ¿Qué cambió?

- Cambio 1
- Cambio 2
- Cambio 3

## Cómo probar

```bash
# 1. Compilar el proyecto
go build -o komyzi.exe ./cmd/cli

# 2. Verificar que compila correctamente
.\komyzi.exe --version

# 3. Probar funcionalidad específica (agregar según el PR):
# .\komyzi.exe [comando] [args]
```

## Checklist

- [ ] Código compila sin errores (`go build ./...`)
- [ ] Tests pasan (`go test ./...`) - si aplica
- [ ] README.md actualizado (si hay cambios de comandos)
- [ ] CHANGELOG.md actualizado (si aplica)
- [ ] Commits siguen [Conventional Commits](https://www.conventionalcommits.org/)
- [ ] No hay comandos "fantasma" en README (documentados pero no implementados)

## Screenshots / Output

Si aplica, agregar output del CLI funcionando:

```
# Ejemplo de output
```

## Notas adicionales

- ¿Hay breaking changes?
- ¿Requiere actualizar dependencias?
- ¿Algo que los reviewers deban saber?

## Referencias

- Issue relacionado: #
- Documentación: #
