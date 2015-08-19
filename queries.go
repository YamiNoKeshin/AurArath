package connection
import (
	"github.com/joernweissenborn/eventual2go"
	"github.com/hashicorp/serf/serf"
)

const (
	EXPORTING = "EXPORTING"
	IMPORTING = "IMPORTING"
)

func IsExporting(d eventual2go.Data) bool {
	return d.(*serf.Query).Name == EXPORTING
}

func IsImporting(d eventual2go.Data) bool {
	return d.(*serf.Query).Name == IMPORTING
}