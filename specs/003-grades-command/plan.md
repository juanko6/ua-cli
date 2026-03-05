# Technical Plan: `ua grades` Implementation

**Feature**: 003-grades-command  
**Status**: Technical Design Complete  
**Est. Effort**: 8-10 hours  
**Dependencies**: Existing auth infrastructure, uacloud adapter

## 🎯 System Architecture Integration

### Ports & Adapters Mapping

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CLI Layer     │    │  Application    │    │   Domain        │
│                 │    │     Layer       │    │     Layer       │
│ ┌─────────────┐ │    │ ┌─────────────┐ │    │ ┌─────────────┐ │
│ │ cmd/grades.go│ │    │ │service/     │ │    │ │domain/      │ │
│ │             │ │    │ │grades/      │ │    │ │grades/       │ │
│ └─────────────┘ │    │ │service.go   │ │    │ │entity.go    │ │
│                 │    │ └─────────────┘ │    │ └─────────────┘ │
│ ┌─────────────┐ │    │                 │    │                 │
│ │ presenter/  │ │    │ ┌─────────────┐ │    │ ┌─────────────┐ │
│ │grades/      │ │────│ │adapter/     │ │────│ │repository/  │ │
│ │text.go      │ │    │ │uacloud/     │ │    │ │grades.go    │ │
│ │json.go      │ │    │ │grades.go    │ │    │ └─────────────┘ │
│ └─────────────┘ │    │ └─────────────┘ │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### New Components Required

#### 1. Domain Layer (`internal/domain/grades/`)
```go
// entity.go
type Grade struct {
    SubjectID      string
    SubjectName    string
    CurrentGrade   string
    Average        float64
    AssessmentCount int
    Status         GradeStatus
    LastUpdated    time.Time
}

type GradeStatus string
const (
    StatusApproved    GradeStatus = "approved"
    StatusPending     GradeStatus = "pending"
    StatusNeedsAttention GradeStatus = "attention"
)

// repository.go
type GradesRepository interface {
    FetchGrades(ctx context.Context, session CookieStore) ([]Grade, error)
    GetLastCheck() time.Time
    SetLastCheck(checkTime time.Time) error
    DetectNewGrades(current []Grade) ([]Grade, error)
}
```

#### 2. Application Layer (`internal/service/grades/`)
```go
// service.go
type GradesService struct {
    repo    GradesRepository
    presenter GradesPresenter
}

func (s *GradesService) GetGrades(ctx context.Context, session CookieStore, opts GradesOptions) (*GradesResult, error)
func (s *GradesService) DetectNewGrades(ctx context.Context, session CookieStore) ([]Grade, error)
```

#### 3. Adapters Layer (`internal/adapters/`)

**uacloud adapter:**
```go
// uacloud/grades.go
type UACloudGradesAdapter struct {
    httpClient *http.Client
    baseUrl   string
}

func (a *UACloudGradesAdapter) FetchGrades(ctx context.Context, cookies []*http.Cookie) ([]Grade, error)
```

**Presenter adapters:**
```go
// presenter/grades/text.go
type TextGradesPresenter struct{}

func (p *TextGradesPresenter) Present(result *GradesResult) string

// presenter/grades/json.go  
type JSONGradesPresenter struct{}

func (p *JSONGradesPresenter) Present(result *GradesResult) string
```

#### 4. CLI Layer (`cmd/ua-cli/`)
```go
// grades.go
var gradesCmd = &cobra.Command{
    Use:   "grades",
    Short: "Show your current grades across all subjects",
    RunE:  runGradesCmd,
    Flags: []pflag.Flag{
        {Name: "json", Shorthand: "j", Usage: "Output in JSON format"},
        {Name: "approved", Usage: "Show only approved subjects"},
        {Name: "pending", Usage: "Show only pending subjects"},
    },
}
```

## 🔧 Technical Implementation Details

### Grade Parsing Strategy

UACloud likely uses HTML tables or JSON APIs for grades. We'll implement a resilient parser:

```go
// internal/adapters/uacloud/grades_parser.go
type GradesParser struct{}

func (p *GradesParser) ParseHTML(html string) ([]Grade, error)
func (p *GradesParser) ParseJSON(jsonData []byte) ([]Grade, error)
```

### Change Detection Algorithm

```go
// internal/domain/grades/tracker.go
type GradeTracker struct {
    storage Storage
}

func (t *GradeTracker) DetectNewGrades(current, previous []Grade) ([]Grade, error) {
    var newGrades []Grade
    
    for _, currentGrade := range current {
        found := false
        for _, prevGrade := range previous {
            if currentGrade.SubjectID == prevGrade.SubjectID {
                found = true
                if currentGrade.LastUpdated.After(prevGrade.LastUpdated) {
                    newGrades = append(newGrades, currentGrade)
                }
                break
            }
        }
        if !found {
            newGrades = append(newGrades, currentGrade)
        }
    }
    
    return newGrades, nil
}
```

### Error Handling Strategy

- **Network errors**: Retry with exponential backoff
- **Parsing errors**: Fallback to error message with raw data for debugging
- **Auth errors**: Redirect to login with helpful message
- **Empty responses**: Graceful handling with informative messages

## 🧪 Testing Strategy

### Unit Tests
- Grade entity validation
- Parser logic for different grade formats
- Change detection algorithm
- Presenter formatting

### Integration Tests  
- UACloud API mocking
- End-to-end command testing
- Session persistence validation

### Acceptance Tests
- `ua grades` command output verification
- `ua grades --json` format validation
- New grade detection flow

## 📊 Performance Considerations

- **Caching**: Cache grades for 1 hour to reduce UACloud calls
- **Parallel loading**: Fetch subject details concurrently
- **Memory optimization**: Stream processing for large datasets
- **Timeout**: 30-second timeout for grade fetching

## 🚨 Risks and Mitigations

### Risk 1: UACloud Grade Format Changes
- **Mitigation**: Abstract parser with configurable patterns
- **Fallback**: Log raw HTML/JSON for debugging

### Risk 2: Performance Issues with Many Subjects  
- **Mitigation**: Implement pagination and lazy loading
- **Monitoring**: Add performance metrics

### Risk 3: Authentication Expiry During Grades Fetch
- **Mitigation**: Refresh session automatically
- **Fallback**: Clear error message asking for re-auth

## 🔗 Dependencies

- **Existing**: auth service, uacloud adapter, presenters
- **New**: None (all dependencies already in project)
- **External**: None (pure Go implementation)

## 📅 Implementation Timeline

1. **Domain Layer**: 2 hours
2. **Application Layer**: 2 hours  
3. **Adapters Layer**: 3 hours
4. **CLI Layer**: 1 hour
5. **Testing**: 2 hours

**Total: 10 hours**