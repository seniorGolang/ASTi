// @asti name="ComplexImports" version=2.0
package compleximports

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	// Локальные импорты с алиасами
	"github.com/seniorGolang/asti/parser/models"
)

// @asti name="ComplexImportService" timeout=60
type ComplexImportService interface {
	// @asti method=ProcessWithSQL retry=5
	ProcessWithSQL(ctx context.Context, query string) (result sql.Result, err error)

	// @asti method=ProcessWithAliasSQL timeout=30
	ProcessWithAliasSQL(ctx context.Context, query string) (result sql.Result, err error)

	// @asti method=ProcessWithDBAlias retry=3
	ProcessWithDBAlias(ctx context.Context, query string) (result sql.Result, err error)

	// @asti method=ProcessJSON timeout=10
	ProcessJSON(ctx context.Context, data interface{}) (result json.RawMessage, err error)

	// @asti method=ProcessJSONAlias retry=2
	ProcessJSONAlias(ctx context.Context, data interface{}) (result json.RawMessage, err error)

	// @asti method=ProcessHTTP2 timeout=15
	ProcessHTTP2(ctx context.Context, stream interface{}) (result interface{}, err error)

	// @asti method=ProcessModels timeout=20
	ProcessModels(ctx context.Context, data models.Annotations) (result models.Annotations, err error)

	// @asti method=ProcessTime timeout=5
	ProcessTime(ctx context.Context, t time.Time) (result time.Time, err error)
}
