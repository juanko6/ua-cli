# Implementation Tasks: `ua grades` Command

**Feature**: 003-grades-command  
**Total Tasks**: 12 atomic tasks  
**Estimated Time**: 8-10 hours  
**Dependencies**: None (atomic sequence)

## 📋 Task Breakdown

### Phase 1: Domain Layer (Tasks 1-3)

#### Task 1: Create Grade Entity
- **File**: `internal/domain/grades/entity.go`
- **Description**: Define core Grade entity and types
- **Acceptance Criteria**:
  - [ ] Grade struct with SubjectID, SubjectName, CurrentGrade, Average, AssessmentCount, Status, LastUpdated
  - [ ] GradeStatus constants (approved, pending, attention)
  - [ ] Validation methods for grade data
  - [ ] Unit tests covering entity validation
- **Estimated Time**: 1 hour
- **Dependencies**: None

#### Task 2: Implement Grade Repository Interface  
- **File**: `internal/domain/grades/repository.go`
- **Description**: Define repository contract for grade operations
- **Acceptance Criteria**:
  - [ ] GradesRepository interface with FetchGrades, GetLastCheck, SetLastCheck, DetectNewGrades methods
  - [ ] InMemoryGradesRepository implementation for testing
  - [ ] Unit tests for repository interface compliance
  - [ ] Error handling for storage operations
- **Estimated Time**: 1.5 hours
- **Dependencies**: Task 1

#### Task 3: Create Grade Tracker
- **File**: `internal/domain/grades/tracker.go`
- **Description**: Implement change detection logic
- **Acceptance Criteria**:
  - [ ] GradeTracker struct with storage dependency
  - [ ] DetectNewGrades method comparing current vs previous grades
  - [ ] Timestamp comparison logic
  - [ ] Unit tests for change detection scenarios
- **Estimated Time**: 1.5 hours
- **Dependencies**: Task 1, Task 2

### Phase 2: Application Layer (Tasks 4-5)

#### Task 4: Create Grades Service
- **File**: `internal/service/grades/service.go`
- **Description**: Implement application logic for grades
- **Acceptance Criteria**:
  - [ ] GradesService struct with repository and presenter dependencies
  - [ ] GetGrades method with filtering options
  - [ ] DetectNewGrades method using tracker
  - [ ] Business logic for grade status calculation
  - [ ] Unit tests for service methods
- **Estimated Time**: 2 hours
- **Dependencies**: Task 1, Task 2, Task 3

#### Task 5: Create Grades Options and Result Types
- **File**: `internal/service/grades/types.go`
- **Description**: Define service contracts and DTOs
- **Acceptance Criteria**:
  - [ ] GradesOptions struct (flags for filtering)
  - [ ] GradesResult struct (output data structure)
  - [ ] NewGrade struct (for change detection)
  - [ ] Unit tests for type validation
- **Estimated Time**: 0.5 hours
- **Dependencies**: Task 4

### Phase 3: Adapters Layer (Tasks 6-8)

#### Task 6: Implement UACloud Grades Adapter
- **File**: `internal/adapters/uacloud/grades.go`
- **Description**: Fetch grades from UACloud API
- **Acceptance Criteria**:
  - [ ] UACloudGradesAdapter struct with HTTP client
  - [ ] FetchGrades method using existing cookie auth
  - [ ] HTML/JSON parsing logic
  - [ ] Error handling for API failures
  - [ ] Unit tests with mocked HTTP responses
- **Estimated Time**: 2 hours
- **Dependencies**: Task 1, existing uacloud adapter

#### Task 7: Create Grades Presenters
- **File**: `internal/adapters/presenter/grades/text.go`
- **File**: `internal/adapters/presenter/grades/json.go`
- **Description**: Implement output formatting for grades
- **Acceptance Criteria**:
  - [ ] TextGradesPresenter with table formatting
  - [ ] JSONGradesPresenter with structured output
  - [ ] Support for highlighting new grades
  - [ ] Unit tests for presenter output
- **Estimated Time**: 1.5 hours
- **Dependencies**: Task 5

#### Task 8: Integrate Repository Implementation
- **File**: `internal/adapters/repo/grades.go`
- **Description**: Concrete implementation of GradesRepository
- **Acceptance Criteria**:
  - [ ] UACloudGradesRepository implementing repository interface
  - [ ] Session management and auth integration
  - [ ] Local caching implementation
  - [ ] Unit tests for repository implementation
- **Estimated Time**: 1 hour
- **Dependencies**: Task 2, Task 6

### Phase 4: CLI Layer (Tasks 9-10)

#### Task 9: Create Grades CLI Command
- **File**: `cmd/ua-cli/grades.go`
- **Description**: Implement CLI command interface
- **Acceptance Criteria**:
  - [ ] Cobra command definition with flags
  - [ ] Command runner with service integration
  - [ ] Help text and usage examples
  - [ ] Flag validation and parsing
  - [ ] Unit tests for command structure
- **Estimated Time**: 1 hour
- **Dependencies**: Task 4, Task 7, Task 8

#### Task 10: Integrate with Root Command
- **File**: `cmd/ua-cli/root.go` (update)
- **Description**: Add grades command to CLI structure
- **Acceptance Criteria**:
  - [ ] Add grades command to root command
  - [ ] Update help documentation
  - [ ] Integration testing with existing commands
- **Estimated Time**: 0.5 hours
- **Dependencies**: Task 9

### Phase 5: Testing and Integration (Tasks 11-12)

#### Task 11: Integration Testing
- **File**: `cmd/ua-cli/grades_test.go`
- **Description**: End-to-end testing of grades command
- **Acceptance Criteria**:
  - [ ] Test with mock UACloud responses
  - [ ] Test authentication integration
  - [ ] Test JSON output format
  - [ ] Test error handling scenarios
- **Estimated Time**: 1 hour
- **Dependencies**: Task 9, Task 10

#### Task 12: Documentation and Cleanup
- **Files**: Update README.md, add examples
- **Description**: Final documentation and code review
- **Acceptance Criteria**:
  - [ ] Update README with new command
  - [ ] Add usage examples
  - [ ] Code review and refactoring
  - [ ] Performance optimization
- **Estimated Time**: 0.5 hours
- **Dependencies**: All previous tasks

## 🔗 Task Dependencies

```
Task 1 → Task 2 → Task 3 → Task 4 → Task 5 → Task 9 → Task 10
    ↓           ↓           ↓           ↓           ↓
Task 6 → Task 8 → Task 7 → Task 11 → Task 12
```

## 🚨 Quality Gates

Each task must meet:
- [ ] All acceptance criteria satisfied
- [ ] Unit tests passing (coverage >80%)
- [ ] Code follows project conventions
- [ ] Error handling implemented
- [ ] Documentation updated

## 📊 Progress Tracking

| Task | Status | Time Est. | Time Actual | Notes |
|------|--------|-----------|-------------|-------|
| Task 1 | ⏳ | 1h | | |
| Task 2 | ⏳ | 1.5h | | |
| Task 3 | ⏳ | 1.5h | | |
| Task 4 | ⏳ | 2h | | |
| Task 5 | ⏳ | 0.5h | | |
| Task 6 | ⏳ | 2h | | |
| Task 7 | ⏳ | 1.5h | | |
| Task 8 | ⏳ | 1h | | |
| Task 9 | ⏳ | 1h | | |
| Task 10 | ⏳ | 0.5h | | |
| Task 11 | ⏳ | 1h | | |
| Task 12 | ⏳ | 0.5h | | |

**Total Estimated: 10 hours**