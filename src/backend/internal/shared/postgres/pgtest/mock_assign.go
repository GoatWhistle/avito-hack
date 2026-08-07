package pgtest

import (
	"fmt"
	"reflect"
)

func assign(dest, values []any) error {
	if len(values) < len(dest) {
		return fmt.Errorf("pgtest: have %d values for %d destinations", len(values), len(dest))
	}

	for i, d := range dest {
		if err := assignOne(d, values[i]); err != nil {
			return fmt.Errorf("pgtest: destination %d: %w", i, err)
		}
	}

	return nil
}

func assignOne(dest, value any) error {
	target := reflect.ValueOf(dest)
	if target.Kind() != reflect.Pointer || target.IsNil() {
		return fmt.Errorf("destination is not a non-nil pointer (%T)", dest)
	}

	elem := target.Elem()

	if value == nil {
		elem.Set(reflect.Zero(elem.Type()))

		return nil
	}

	src := reflect.ValueOf(value)

	switch {
	case src.Type().AssignableTo(elem.Type()):
		elem.Set(src)
	case src.Type().ConvertibleTo(elem.Type()):
		elem.Set(src.Convert(elem.Type()))
	default:
		return fmt.Errorf("cannot assign %T to %s", value, elem.Type())
	}

	return nil
}
