package pipelines

import "reflect"

func notifySubscriber(fn interface{}, values []interface{}) []interface{} {
	refFn := reflect.TypeOf(fn)
	fnArgs := make([]reflect.Value, 0, refFn.NumIn())
	interfaceVals := make([]interface{}, 0)

	for argIndex := 0; argIndex < refFn.NumIn(); argIndex++ {
		if len(values) == argIndex {
			return interfaceVals
		}

		providedVal := values[argIndex]

		// Variadic arguments need special treatment
		if refFn.IsVariadic() && refFn.In(argIndex).Kind() == reflect.Slice && argIndex == refFn.NumIn()-1 {
			sliceType := refFn.In(argIndex).Elem()

			for _, innerVal := range values[argIndex:len(values)] {
				if innerVal == nil {
					fnArgs = append(fnArgs, reflect.New(sliceType).Elem())
					continue
				}

				if !reflect.TypeOf(innerVal).AssignableTo(sliceType) {
					// Slice does not match received data, skipping this subscriber
					return interfaceVals
				}
				fnArgs = append(fnArgs, reflect.ValueOf(innerVal))
			}
			// Finish loop as we have filled in all data to the slice
			break
		} else {
			argType := refFn.In(argIndex)
			if providedVal == nil {
				values[argIndex] = reflect.New(argType).Elem()
				providedVal = values[argIndex]
			}

			if !reflect.TypeOf(providedVal).AssignableTo(argType) {
				// Method signature not compatible with this input. Skipping subscriber
				return interfaceVals
			}

			fnArgs = append(fnArgs, reflect.ValueOf(values[argIndex]))

			if argIndex == refFn.NumIn()-1 {
				if refFn.NumIn() != len(fnArgs) {
					// Skipping non-slice overflow
					return interfaceVals
				}
			}
		}

	}

	returnVals := reflect.ValueOf(fn).Call(fnArgs)

	for _, val := range returnVals {
		interfaceVals = append(interfaceVals, val.Interface())
	}
	return interfaceVals
}
