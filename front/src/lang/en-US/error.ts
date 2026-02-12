export default {
  // Auth (41000-41099)
  auth: {
    41005: 'Invalid or expired token',
    66001: 'Authentication failed',
    66002: 'Invalid username or password',
    66003: 'Failed to generate token',
    66004: 'Invalid token',
    66005: 'Token expired',
  },

  // Server (60000-60059)
  server: {
    60001: 'Server not found',
    60002: 'Server alias already exists',
    60003: 'Server not found by UUID',
    60004: 'Server already exists',
    60005: 'Server update failed',
    60006: 'Server delete failed',
    60021: 'Invalid parameter',
    60022: 'Invalid auth type',
    60023: 'Invalid SSH configuration',
    60024: 'SSH connection test failed',
    60041: 'Connect failed',
    60042: 'Agent is offline',
    60043: 'Agent request failed',
  },

  // Config (65000-65019)
  config: {
    65001: 'Config not found',
    65002: 'Config key already exists',
    65003: 'Invalid config key',
    65004: 'Invalid config value',
    65005: 'Config update failed',
    65006: 'Config delete failed',
  },

  // Application (71000-71019)
  application: {
    71001: 'Application not found',
    71002: 'Application already exists',
    71003: 'Invalid application name',
    71004: 'Invalid application type',
    71005: 'Invalid application configuration',
    71006: 'Application update failed',
    71007: 'Application delete failed',
  },

  // Deployment (72000-72039)
  deployment: {
    72001: 'Deployment not found',
    72002: 'Application already deployed to this server',
    72003: 'Application not deployed to this server',
    72004: 'Failed to generate deploy ID',
    72005: 'Failed to create deployment record',
    72006: 'Invalid deployment configuration',
    72007: 'Docker-compose container name conflict detected',
    72008: 'Docker-compose port conflict detected',
    72009: 'Docker-compose volume conflict detected',
    72010: 'Docker-compose network conflict detected',
    72021: 'Failed to send request to agent',
    72022: 'Failed to parse agent response',
    72023: 'Agent deployment failed',
    72024: 'Agent delete application failed',
    72025: 'Agent stop application failed',
    72026: 'Agent start application failed',
    72027: 'Agent operation failed',
  },

  // AppStore (73000-73019)
  appStore: {
    73001: 'Application store not found',
    73002: 'Application store already exists',
    73003: 'Invalid compose content',
    73004: 'Unsupported application type',
    73005: 'Invalid app store configuration',
    73006: 'App store update failed',
    73007: 'App store delete failed',
  },

  // Script (80000-80039)
  script: {
    80001: 'Script not found',
    80002: 'Script already exists',
    80003: 'Invalid script content',
    80004: 'Unsupported script type',
    80005: 'Script not deployed',
    80021: 'Script execution failed',
    80022: 'Script execution timeout',
    80023: 'Server not found',
  },

  // Monitor (81000-81019)
  monitor: {
    81001: 'Monitor request failed',
    81002: 'Invalid monitor configuration',
    81003: 'Monitor data not found',
    81004: 'Server not found',
  },

  // Common errors
  common: {
    networkError: 'Network error, please try again later',
    unknownError: 'Unknown error',
    requestError: 'Request failed',
  },
}
