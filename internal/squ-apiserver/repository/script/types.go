package script

const (
	// 脚本类型
	TypeInstall = "install"   // 安装脚本
	TypeUpdate  = "update"    // 更新脚本
	TypeUninstall = "uninstall" // 卸载脚本
	TypeCheck   = "check"     // 检查脚本
	TypeBackup  = "backup"    // 备份脚本
	TypeRestore = "restore"   // 恢复脚本

	// 脚本分类
	CategoryDocker      = "docker"      // Docker 相关
	CategoryDatabase    = "database"    // 数据库相关
	CategorySystem      = "system"      // 系统相关
	CategoryMonitoring  = "monitoring"  // 监控相关
	CategoryStorage     = "storage"     // 存储相关
	CategoryNetwork     = "network"     // 网络相关
	CategorySecurity    = "security"    // 安全相关

	// 脚本状态
	StatusActive     = "active"     // 活跃
	StatusDeprecated = "deprecated" // 已废弃
	StatusTesting    = "testing"    // 测试中

	// 脚本执行状态
	ExecutionStatusPending   = "pending"   // 待执行
	ExecutionStatusRunning   = "running"   // 执行中
	ExecutionStatusSuccess   = "success"   // 成功
	ExecutionStatusFailed    = "failed"    // 失败
	ExecutionStatusCancelled = "cancelled" // 已取消
)
