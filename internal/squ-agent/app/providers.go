package app

// 这个文件用于放应用级 provider / assembler。
// 当模块越来越多时，可以把数据库、缓存、外部客户端的装配逻辑逐步抽到这里，
// 避免 options.go 或 app.go 变得过胖。
