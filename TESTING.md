# Test Coverage Improvement Plan

## Current Status (2026-03-04)

### Test Coverage Summary
- **Overall**: 19.0%
- **Core Modules**: ~25% (needs improvement)
- **Router**: ~11% (critical)
- **Tree**: ~10% (critical)
- **Upload**: 0% (no tests)
- **Session**: 23.5%
- **JSON**: 77.8% ✅
- **String**: 78.6% ✅

## Goals

### Phase 1: Core Testing (Target: 35%+)
- [ ] Add router_test.go tests
- [ ] Add tree_test.go tests  
- [ ] Add uploadfile_test.go tests
- [ ] Improve group_test.go coverage

### Phase 2: Edge Cases (Target: 45%+)
- [ ] Route conflict tests
- [ ] Parameter parsing edge cases
- [ ] Session concurrent tests
- [ ] Middleware chain tests

### Phase 3: Benchmarks
- [ ] Router matching benchmarks
- [ ] Session read/write benchmarks
- [ ] Middleware chain benchmarks

## Running Tests

```bash
# Run all tests with coverage
go test ./... -coverprofile=coverage.out

# View coverage report
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
```

## Test Naming Convention

```
Test<FunctionName>_<Scenario>_<ExpectedResult>

Examples:
- TestRouter_AddRoute_ValidPath_Success
- TestRouter_AddRoute_EmptyPath_Error
- TestGroup_Use_MiddlewareChain_Order
```

## CI Integration

GitHub Actions workflow in `.github/workflows/test.yml` runs on:
- Push to aicode, master branches
- Pull requests
