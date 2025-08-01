// @asti name="EdgeCases" version=1.0
package edgecases

import (
	"context"
	"time"
)

// Интерфейс без аннотации - должен быть пропущен парсером
type UnannotatedInterface interface {
	Method(ctx context.Context, param string) (result string, err error)
}

// @asti name="ValidInterface" timeout=30
type ValidInterface interface {
	// @asti method=ValidMethod retry=3
	ValidMethod(ctx context.Context, param string) (result string, err error)
}

// Интерфейс с методом без context.Context - должен быть пропущен
// @asti name="InvalidContextInterface" timeout=30
type InvalidContextInterface interface {
	// @asti method=InvalidMethod retry=3
	InvalidMethod(param string) (result string, err error)
}

// Интерфейс с методом без error - должен быть пропущен
// @asti name="InvalidErrorInterface" timeout=30
type InvalidErrorInterface interface {
	// @asti method=InvalidMethod retry=3
	InvalidMethod(ctx context.Context, param string) (result string)
}

// Интерфейс с неправильным именем первого параметра - должен быть пропущен
// @asti name="UnnamedResultsInterface" timeout=30
type UnnamedResultsInterface interface {
	// @asti method=InvalidMethod retry=3
	InvalidMethod(ctx context.Context, param string) (string, error)
}

// Интерфейс с очень длинными именами
// @asti name="VeryLongInterfaceNameThatExceedsNormalLimitsAndTestsHowTheParserHandlesExtremelyLongIdentifiers" timeout=60
type VeryLongInterfaceNameThatExceedsNormalLimitsAndTestsHowTheParserHandlesExtremelyLongIdentifiers interface {
	// @asti method=VeryLongMethodNameThatExceedsNormalLimitsAndTestsHowTheParserHandlesExtremelyLongIdentifiers retry=5
	VeryLongMethodNameThatExceedsNormalLimitsAndTestsHowTheParserHandlesExtremelyLongIdentifiers(ctx context.Context, veryLongParameterName string) (veryLongResultName string, err error)
}

// Интерфейс с множественными аннотациями
// @asti name="MultipleAnnotations" version=2.0 timeout=45
// @asti category=service model=business
// @asti author="Senior Developer" team="Backend"
type MultipleAnnotations interface {
	// @asti method=Method1 retry=3 timeout=10
	// @asti validation=strict authorization=required
	Method1(ctx context.Context, param1 string) (result1 string, err error)

	// @asti method=Method2 retry=5 timeout=20
	// @asti validation=partial authorization=user
	Method2(ctx context.Context, param2 int) (result2 int, err error)
}

// Интерфейс с очень сложными типами
// @asti name="ComplexTypesInterface" timeout=90
type ComplexTypesInterface interface {
	// @asti method=ProcessComplexMap retry=3
	ProcessComplexMap(ctx context.Context, data map[string]map[string][]map[string]interface{}) (result map[string]map[string][]map[string]interface{}, err error)

	// @asti method=ProcessComplexSlice timeout=30
	ProcessComplexSlice(ctx context.Context, data [][][][][]interface{}) (result [][][][][]interface{}, err error)

	// @asti method=ProcessComplexFunction timeout=45
	ProcessComplexFunction(ctx context.Context, fn func(context.Context, string, int) func(map[string]interface{}) error) (result func(context.Context, string, int) func(map[string]interface{}) error, err error)
}

// Интерфейс с аннотациями, содержащими специальные символы
// @asti name="SpecialChars" version="1.0" timeout=30
// @asti description="Interface with special characters in annotations"
// @asti tags="test,complex,special"
type SpecialChars interface {
	// @asti method=MethodWithSpecialChars retry=3
	// @asti description="Method with special characters: @#$%^&*()"
	// @asti validation="strict,required" authorization="user,admin"
	MethodWithSpecialChars(ctx context.Context, param string) (result string, err error)
}

// Интерфейс с очень большим количеством методов
// @asti name="ManyMethodsInterface" timeout=120
type ManyMethodsInterface interface {
	// @asti method=Method1 retry=1
	Method1(ctx context.Context, param string) (result string, err error)

	// @asti method=Method2 retry=2
	Method2(ctx context.Context, param int) (result int, err error)

	// @asti method=Method3 retry=3
	Method3(ctx context.Context, param bool) (result bool, err error)

	// @asti method=Method4 retry=4
	Method4(ctx context.Context, param float64) (result float64, err error)

	// @asti method=Method5 retry=5
	Method5(ctx context.Context, param time.Time) (result time.Time, err error)

	// @asti method=Method6 retry=6
	Method6(ctx context.Context, param []string) (result []string, err error)

	// @asti method=Method7 retry=7
	Method7(ctx context.Context, param map[string]interface{}) (result map[string]interface{}, err error)

	// @asti method=Method8 retry=8
	Method8(ctx context.Context, param interface{}) (result interface{}, err error)

	// @asti method=Method9 retry=9
	Method9(ctx context.Context, param *string) (result *string, err error)

	// @asti method=Method10 retry=10
	Method10(ctx context.Context, param chan string) (result chan string, err error)
}
