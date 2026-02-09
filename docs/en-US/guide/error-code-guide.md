# Squirrel Error Code Specification

## Overview

This document defines the specification and usage guidelines for all API error codes in the Squirrel system. A unified error code design helps with:
- Quickly identifying problem sources
- Distinguishing errors from different modules
- Facilitating debugging and log analysis
- Providing user-friendly error messages

## Error Code Design Principles

### Naming Conventions
- **Uniform 5-digit format**: All error codes are 5 digits for easy memory and recognition
- **Modular segmentation**: Each module occupies an independent space of 1000 error codes
- **Functional subdivision**: Within each module, error codes are grouped in sets of 20
- **Reserved expansion**: Each module reserves sufficient space for future growth

### Error Code Structure
```
[Module ID][Functional Group][Specific Error]
  60      01        01
  ↓        ↓        ↓
  Module  Function   Specific Error
```

## Base Error Codes

### 0xxxx: General Status
| Error Code | Constant Name | Description | Usage Scenario |
|------------|---------------|-------------|----------------|
| 0 | `CodeSuccess` | Operation successful | All successful API responses |

### 41xxx: Common Errors
| Error Code | Constant Name | Description | Usage Scenario |
|------------|---------------|-------------|----------------|
| 41001 | `ErrCodeParameter` | Parameter error | Request parameter validation failed |
| 41002 | `ErrUserOrPassword` | Username or password error | Authentication failed |

### 50xxx: Database Errors
| Error Code | Constant Name | Description | Usage Scenario |
|------------|---------------|-------------|----------------|
| 50000 | `ErrSQL` | SQL error | Database operation exception |
| 50001 | `ErrSQLNotFound` | Data not found | Database query returned no results |
| 50002 | `ErrSQLNotUnique` | Data not unique | Unique constraint conflict |
| 50003 | `ErrDuplicatedKey` | Duplicated key | Primary key or unique key conflict |

**Definition Location**: `internal/pkg/response/common.go`

---

## Module Error Codes

### 60xxx: Server Management

**Definition Location**: `internal/squ-apiserver/handler/server/res/response_code.go`

#### 60000-60019: Basic Operations (CRUD)
| Error Code | Constant Name | Description | Usage Scenario |
|------------|---------------|-------------|----------------|
| 60001 | `ErrServerNotFound` | Server not found | No server found based on query conditions |
| 60002 | `ErrServerAliasExists` | Server alias already exists | Alias conflict during add/update |
| 60003 | `ErrServerUUIDNotFound` | Server not found by UUID | UUID does not exist during agent registration |
| 60004 | `ErrServerAlreadyExists` | Server already exists | Duplicate addition of the same server |
| 60005 | `ErrServerUpdateFailed` | Server update failed | Update operation exception |
| 60006 | `ErrServerDeleteFailed` | Server delete failed | Delete operation exception |

#### 60020-60039: Validation
| Error Code | Constant Name | Description | Usage Scenario |
|------------|---------------|-------------|----------------|
| 60021 | `ErrInvalidParameter` | Parameter validation failed | Request parameters do not meet requirements |
| 60022 | `ErrInvalidAuthType` | Invalid authentication type | Authentication type is not password/key/password_key |
| 60023 | `ErrInvalidSSHConfig` | Invalid SSH configuration | SSH configuration format or content error |
| 60024 | `ErrSSHTestFailed` | SSH connection test failed | Unable to establish SSH connection |

#### 60040-60059: Agent Communication
| Error Code | Constant Name | Description | Usage Scenario |
|------------|---------------|-------------|----------------|
| 60041 | `ErrConnectFailed` | Connection failed | Network connection exception |
| 60042 | `ErrAgentOffline` | Agent offline | Agent service unavailable |
| 60043 | `ErrAgentRequestFailed` | Agent request failed | Failed to send request to agent |

---

### 70xxx: App Store

**Definition Location**: `internal/squ-apiserver/handler/app_store/res/response_code.go`

| Error Code | Constant Name | Description | Usage Scenario |
|------------|---------------|-------------|----------------|
| 70001 | `ErrAppStoreNotFound` | App store not found | Queried app store does not exist |
| 70002 | `ErrDuplicateAppStore` | App store already exists | Duplicate app store addition |
| 70003 | `ErrInvalidComposeContent` | Invalid Compose content | Docker Compose format error |
| 70004 | `ErrUnsupportedAppType` | Unsupported application type | Application type not in supported list |

---

### 71xxx: Application

**Definition Location**: `internal/squ-apiserver/handler/application/res/response_code.go`

| Error Code | Constant Name | Description | Usage Scenario |
|------------|---------------|-------------|----------------|
| 71001 | `ErrApplicationNotFound` | Application not found | Queried application does not exist |
| 71002 | `ErrDuplicateApplication` | Application already exists | Duplicate application addition |
| 71003 | `ErrInvalidAppConfig` | Invalid application configuration | Application configuration format or content error |

---

### 72xxx: Deployment

**Definition Location**: `internal/squ-apiserver/handler/deployment/res/response_code.go`

#### 72000-72019: Basic Operations
| Error Code | Constant Name | Description | Usage Scenario |
|------------|---------------|-------------|----------------|
| 72001 | `ErrDeploymentNotFound` | Deployment record not found | Queried deployment does not exist |
| 72002 | `ErrAlreadyDeployed` | Application already deployed to this server | Duplicate deployment |
| 72003 | `ErrApplicationNotDeployed` | Application not deployed to this server | Operating on an undeployed application |
| 72004 | `ErrDeployIDGenerateFailed` | Failed to generate deployment ID | ID generator exception |
| 72005 | `ErrCreateDeploymentRecordFailed` | Failed to create deployment record | Database write failed |

#### 72020-72039: Agent Related
| Error Code | Constant Name | Description | Usage Scenario |
|------------|---------------|-------------|----------------|
| 72021 | `ErrAgentRequestFailed` | Agent request failed | Failed to send request to agent |
| 72022 | `ErrAgentResponseParseFailed` | Agent response parse failed | Response data format error |
| 72023 | `ErrAgentDeployFailed` | Agent deployment failed | Agent-side deployment operation failed |
| 72024 | `ErrAgentDeleteFailed` | Agent delete application failed | Agent-side delete operation failed |
| 72025 | `ErrAgentStopFailed` | Agent stop application failed | Agent-side stop operation failed |
| 72026 | `ErrAgentStartFailed` | Agent start application failed | Agent-side start operation failed |
| 72027 | `ErrAgentOperationFailed` | Agent operation failed | Agent-side other operation failed |

---

### 80xxx: Script

**Definition Location**: `internal/squ-apiserver/handler/script/res/response_code.go`

#### 80000-80019: Basic Operations
| Error Code | Constant Name | Description | Usage Scenario |
|------------|---------------|-------------|----------------|
| 80001 | `ErrScriptNotFound` | Script not found | Queried script does not exist |
| 80002 | `ErrDuplicateScript` | Script already exists | Duplicate script addition |
| 80003 | `ErrInvalidScriptContent` | Invalid script content | Script content format error |
| 80004 | `ErrUnsupportedScriptType` | Unsupported script type | Script type not in supported list |
| 80005 | `ErrScriptNotDeployed` | Script not deployed | Operating on an undeployed script |

#### 80020-80039: Execution
| Error Code | Constant Name | Description | Usage Scenario |
|------------|---------------|-------------|----------------|
| 80021 | `ErrScriptExecutionFailed` | Script execution failed | Script runtime exception |
| 80022 | `ErrScriptTimeout` | Script execution timeout | Script execution exceeded time limit |
| 80023 | `ErrServerNotFound` | Server not found (script related) | Server associated with script does not exist |

---

### 81xxx: Monitor

**Definition Location**: `internal/squ-apiserver/handler/monitor/res/response_code.go`

| Error Code | Constant Name | Description | Usage Scenario |
|------------|---------------|-------------|----------------|
| 81001 | `ErrMonitorFailed` | Monitor request failed | Failed to get monitoring data |
| 81002 | `ErrInvalidMonitorConfig` | Invalid monitoring configuration | Monitoring configuration format error |
| 81003 | `ErrMonitorDataNotFound` | Monitoring data not found | Monitoring data does not exist |
| 81004 | `ErrServerNotFound` | Server not found (monitor related) | Server associated with monitoring does not exist |

---

## Error Code Allocation Overview

### Allocation Rules

| Module | Error Code Range | Description |
|--------|------------------|-------------|
| General Status | 0 | Success status |
| Common Errors | 41000-41999 | Common errors like parameters, authentication |
| Database Errors | 50000-50999 | SQL operation related errors |
| Server Management | 60000-60999 | Server module |
| App Store | 70000-70999 | App Store module |
| Application | 71000-71999 | Application module |
| Deployment | 72000-72999 | Deployment module |
| Script | 80000-80999 | Script module |
| Monitor | 81000-81999 | Monitor module |
| Reserved Expansion | 82000-99999 | Future functional modules |

### Extension Guidelines

When adding new error codes, please follow these steps:

1. **Confirm module range**: Check the error code allocation table for the corresponding module
2. **Select functional group**: Find the appropriate functional group within the module (20 codes per group)
3. **Define constant**: Define the constant in the corresponding `response_code.go`
4. **Register error code**: Register the error code and description in the `RegisterCode()` function
5. **Update documentation**: Update this document accordingly

**Example**:
```go
// Define constant
const (
    ErrNewFeatureFailed = 60045 // New feature failed
)

// Register error code
func RegisterCode() {
    response.Register(ErrNewFeatureFailed, "new feature failed")
}
```

---

## Error Code Usage Examples

### Basic Error Return

```go
func (s *Server) Get(id uint) response.Response {
    daoS, err := s.Repository.Get(id)
    if err != nil {
        return response.Error(res.ErrServerNotFound)
    }
    return response.Success(daoS)
}
```

### Database Error Handling

```go
func (s *Server) Add(request req.Server) response.Response {
    modelReq := s.requestToModel(request)
    
    err := s.Repository.Add(&modelReq)
    if err != nil {
        // Check if it's a duplicate key error
        if errors.Is(err, gorm.ErrDuplicatedKey) {
            return response.Error(res.ErrServerAlreadyExists)
        }
        return response.Error(model.ReturnErrCode(err))
    }
    
    return response.Success("success")
}
```

### Parameter Validation Error

```go
func (s *Server) Update(request req.Server) response.Response {
    if request.Hostname == "" {
        return response.Error(res.ErrInvalidParameter)
    }
    
    modelReq := s.requestToModel(request)
    modelReq.ID = request.ID
    
    err := s.Repository.Update(&modelReq)
    if err != nil {
        return response.Error(res.ErrServerUpdateFailed)
    }
    
    return response.Success("success")
}
```

---

## Important Notes

### 1. Error Code Uniqueness
- Ensure error codes are unique throughout the system
- Different modules must not share the same error code
- Check for existence before adding new error codes

### 2. Error Message Friendliness
- Error messages should be clear, accurate, and easy to understand
- Use English descriptions for internationalization
- Avoid exposing sensitive information

### 3. Logging
- Return error codes while logging detailed information
- Logs should include error context information
- Facilitate problem tracking and debugging

### 4. Backward Compatibility
- Published error codes must not be arbitrarily changed
- New error codes should be version-marked in the documentation
- Deprecated error codes should be retained for at least one major version

---

## Appendix: Related Files

### Core Files
- `internal/pkg/response/common.go` - Base error code definitions
- `internal/pkg/response/response.go` - Response structure and error handling functions

### Module Files
- `internal/squ-apiserver/handler/server/res/response_code.go` - Server module
- `internal/squ-apiserver/handler/app_store/res/response_code.go` - App Store module
- `internal/squ-apiserver/handler/application/res/response_code.go` - Application module
- `internal/squ-apiserver/handler/deployment/res/response_code.go` - Deployment module
- `internal/squ-apiserver/handler/script/res/response_code.go` - Script module
- `internal/squ-apiserver/handler/monitor/res/response_code.go` - Monitor module

---

## Version History

| Version | Date | Description |
|---------|------|-------------|
| v1.0 | 2026-02-08 | Initial version, unified error code specification |

---

**Maintainer**: Squirrel Development Team
**Last Updated**: 2026-02-08
