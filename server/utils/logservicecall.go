package utils

import (
	"fmt"
	"time"
)

func LogServiceCall(serviceName, methodName string, now time.Time) {
	elapsed := time.Since(now)
	fmt.Printf("Service: %s, Method: %s, Duration: %s\n", serviceName, methodName, elapsed)
}
