Test Coverage Report


Summary

Command run: `go test ./... -coverprofile=coverage.out -covermode=atomic`
Goal: raise package coverage to >= 90%

Per-package coverage

* internal/config: 100.0%
* internal/email: 86.1%
* internal/handlers: 90.5%
* internal/templatecache: 100.0%
* pkg/middleware: 100.0%
* cmd/server: 0.0% (no tests)
* docs: 0.0% (no tests)

Files / packages below 90%

internal/email is at 86.1% and contains code paths not yet covered by unit tests. This is the primary gap to reach an overall repo target of 90%.
cmd/server and docs contain no unit tests; decide whether to include them in the target scope.

Key uncovered areas (internal/email)

* RenderTemplate parsing/execution error branches when template parsing fails or the cache returns unexpected values.
* cleanupExpiredOTPs ticker-driven goroutine (the background branch is not directly covered).
* Extra concurrency/error branches in SendBulkEmail and any untested small helpers.

Recommendations & Next Actions

1) Target internal/email to reach >= 90%:
   - Add tests that force RenderTemplate failures by temporarily registering invalid templates in predefinedTemplates and asserting errors from RenderTemplate and SendTemplateEmail.
   - Add a testable cleanup helper or make the ticker injectable so the cleanup logic can be invoked synchronously in tests to assert OTP removal behavior.
   - Add more SendBulkEmail tests for concurrency/batch size edges and simulated failures.

2) (Optional) Add minimal tests for cmd/server startup (smoke test) to raise coverage if needed.

3) Produce an overall coverage HTML report:

```bash
go test ./... -coverprofile=coverage.out -covermode=atomic
go tool cover -html=coverage.out -o coverage.html
```

Actions I can take next

* Implement the RenderTemplate failure tests and a deterministic cleanupExpiredOTPs test harness (recommended next step — highest ROI for coverage).
* Implement a cmd/server smoke test if you want to include it in the coverage target.

If you want me to proceed, tell me which of the recommended next actions to run and I'll implement them and re-run the full test suite.
