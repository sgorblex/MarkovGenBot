package assert

import (
	"log"
)

// If predicate is false, triggers a panic
func Simple(predicate bool) {
	if !predicate {
		log.Panicln("Assertion error: unspecified")
	}
}

// If predicate is false, triggers a panic with the given message
func String(predicate bool, message string) {
	if !predicate {
		log.Panicln("Assertion error: ", message)
	}
}
