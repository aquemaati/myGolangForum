package config_test

import (
	"os"
	"testing"

	"github.com/aquemaati/myGolangForum.git/internal/config"
)

func TestLoadEnvFileWithTempFile(t *testing.T) {
	// Contenu du fichier .env temporaire
	envContent := `TEST_VARIABLE=TempTestValue
ANOTHER_VARIABLE=TempAnotherValue`

	// Créer un fichier temporaire
	tempFile, err := os.CreateTemp("", "temp-env-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Supprimer le fichier temporaire après le test

	// Écrire le contenu dans le fichier temporaire
	if _, err := tempFile.Write([]byte(envContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Fermer le fichier pour garantir que son contenu est bien écrit
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Appeler la fonction LoadEnvFile en passant le nom du fichier temporaire
	if err := config.LoadEnvFile(tempFile.Name()); err != nil {
		t.Fatalf("Failed to load environment file: %v", err)
	}

	// Vérifier que les variables sont chargées correctement
	expectedVars := map[string]string{
		"TEST_VARIABLE":    "TempTestValue",
		"ANOTHER_VARIABLE": "TempAnotherValue",
	}

	for key, expectedValue := range expectedVars {
		actualValue := os.Getenv(key)
		if actualValue != expectedValue {
			t.Errorf("Expected %s to be %s, got %s", key, expectedValue, actualValue)
		}
	}
}
