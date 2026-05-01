# Test Instructions

- All tests for TuningAccuracy:
  ```
  go test -v ./internal/logic -run TestTuningAccuracy
  ```

- Specific tests for TuningAccuracy:
  ```
   go test -v ./internal/logic -run TestTuningAccuracy/Kepler
   go test -v ./internal/logic -run TestTuningAccuracy/Equal
   go test -v ./internal/logic -run TestTuningAccuracy/Pythagorean
  ```