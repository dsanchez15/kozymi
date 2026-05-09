#!/bin/bash
set -e

# Script de Release para Komyzi
# Uso: ./scripts/release.sh <version>
# Ejemplo: ./scripts/release.sh 0.2.0

VERSION=$1

if [ -z "$VERSION" ]; then
    echo "Error: Debes especificar una versión"
    echo "Uso: ./scripts/release.sh <version>"
    echo "Ejemplo: ./scripts/release.sh 0.2.0"
    exit 1
fi

# Validar formato de versión (SemVer)
if ! echo "$VERSION" | grep -E '^[0-9]+\.[0-9]+\.[0-9]+$' > /dev/null; then
    echo "Error: La versión debe seguir SemVer (ej: 0.2.0)"
    exit 1
fi

TAG="v${VERSION}"

echo "========================================"
echo "  Komyzi Release Script"
echo "========================================"
echo ""
echo "Versión: $VERSION"
echo "Tag: $TAG"
echo ""

# 1. Validar que estamos en main
echo "1. Validando rama..."
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo "❌ Error: Debes estar en la rama 'main' para hacer un release"
    echo "   Rama actual: $CURRENT_BRANCH"
    echo "   Ejecuta: git checkout main"
    exit 1
fi
echo "   ✅ En rama main"

# 2. Validar que no hay cambios sin commitear
echo "2. Validando estado del repositorio..."
if ! git diff-index --quiet HEAD --; then
    echo "❌ Error: Hay cambios sin commitear"
    echo "   Haz commit o stash de los cambios primero"
    git status
    exit 1
fi
echo "   ✅ Repositorio limpio"

# 3. Validar que la versión en código coincide
echo "3. Validando versión en código..."
CODE_VERSION=$(grep 'var version = ' cmd/cli/main.go | sed 's/.*"\(.*\)".*/\1/')
if [ "$CODE_VERSION" != "$VERSION" ]; then
    echo "❌ Error: La versión en cmd/cli/main.go ($CODE_VERSION) no coincide con $VERSION"
    echo "   Actualiza la versión antes de continuar"
    exit 1
fi
echo "   ✅ Versión en código: $CODE_VERSION"

# 4. Validar que existe entrada en CHANGELOG
echo "4. Validando CHANGELOG.md..."
if ! grep -q "\[${VERSION}\]" CHANGELOG.md; then
    echo "❌ Error: No existe entrada para [$VERSION] en CHANGELOG.md"
    echo "   Agrega los cambios antes de continuar"
    exit 1
fi
echo "   ✅ CHANGELOG.md tiene entrada para $VERSION"

# 5. Validar que existe entrada en homebrew formula
echo "5. Validando Homebrew formula..."
if ! grep -q "version \"${VERSION}\"" homebrew/komyzi.rb; then
    echo "❌ Error: La fórmula de Homebrew no tiene la versión $VERSION"
    echo "   Actualiza homebrew/komyzi.rb antes de continuar"
    exit 1
fi
echo "   ✅ Homebrew formula tiene versión $VERSION"

# 6. Validar que el tag no existe
echo "6. Validando tag..."
if git rev-parse "$TAG" >/dev/null 2>&1; then
    echo "❌ Error: El tag $TAG ya existe"
    echo "   Usa una versión diferente"
    exit 1
fi
echo "   ✅ Tag $TAG disponible"

# 7. Compilar localmente para validar
echo "7. Validando compilación..."
if ! go build -o /tmp/komyzi-test ./cmd/cli; then
    echo "❌ Error: Falló la compilación"
    exit 1
fi
rm -f /tmp/komyzi-test
echo "   ✅ Compilación exitosa"

# 8. Ejecutar tests
echo "8. Ejecutando tests..."
if ! go test ./...; then
    echo "❌ Error: Tests fallaron"
    exit 1
fi
echo "   ✅ Tests pasaron"

echo ""
echo "========================================"
echo "  Validaciones completadas ✅"
echo "========================================"
echo ""
echo "Resumen:"
echo "  - Rama: main"
echo "  - Repositorio limpio"
echo "  - Versión en código: $VERSION"
echo "  - CHANGELOG: Actualizado"
echo "  - Homebrew: Actualizado"
echo "  - Tag: $TAG (nuevo)"
echo "  - Compilación: OK"
echo "  - Tests: OK"
echo ""
echo "¿Deseas crear el tag y publicar el release?"
echo "Escribe 'RELEASE' para continuar:"
read CONFIRMATION

if [ "$CONFIRMATION" != "RELEASE" ]; then
    echo "Cancelado por el usuario"
    exit 0
fi

# 9. Crear tag
echo ""
echo "9. Creando tag $TAG..."
git tag -a "$TAG" -m "Release $TAG"
echo "   ✅ Tag creado"

# 10. Push del tag
echo "10. Publicando tag..."
git push origin "$TAG"
echo "   ✅ Tag publicado"

echo ""
echo "========================================"
echo "  🎉 Release $TAG publicado!"
echo "========================================"
echo ""
echo "El workflow de GitHub Actions ahora compilará y creará la release automáticamente."
echo "Puedes ver el progreso en:"
echo "  https://github.com/komyzi/komyzi/actions"
echo ""
echo "La release aparecerá en:"
echo "  https://github.com/komyzi/komyzi/releases"
