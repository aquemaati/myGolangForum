# Nom du projet
APP_NAME := my-go-app

# Variables
MAIN_FILE := main.go

# Spécifiez les outils à utiliser
GO := go

# Par défaut, exécuter `go run main.go`
default: run

# Commande pour exécuter `go run main.go`
run:
	$(GO) run $(MAIN_FILE)

# Commande pour nettoyer les fichiers binaires (si nécessaire)
clean:
	$(GO) clean

# Commande pour installer les dépendances
deps:
	$(GO) mod tidy

# Compilation du projet en un fichier binaire
build:
	$(GO) build -o $(APP_NAME) $(MAIN_FILE)

# Exécution des tests
test:
	$(GO) test -v ./...

# Commande générale
.PHONY: default run clean deps build test
