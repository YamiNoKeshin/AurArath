package aurarath
import (
	"testing"
	"aurarath/config"
)

func TestImportExportedEvents(t *testing.T) {
	i := NewImport(new(AppDescriptor),config.DefaultLocalhost())
	e := NewExport(new(AppDescriptor),config.DefaultLocalhost())
}