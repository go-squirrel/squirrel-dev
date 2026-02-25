// 监控数据工具函数

import type { DiskIORecord, NetworkIORecord } from '@/types/monitor'

// 磁盘IO速率记录（前端计算后）
export interface DiskIOSpeedRecord {
  collect_time: string
  read_speed: number      // bytes/s
  write_speed: number     // bytes/s
}

// 网络IO速率记录（前端计算后）
export interface NetworkIOSpeedRecord {
  collect_time: string
  upload_speed: number    // bytes/s
  download_speed: number  // bytes/s
}

/**
 * 将磁盘IO累计值转换为速率
 * @param records 按时间正序排列的历史记录（旧→新）
 * @returns 速率记录数组
 *
 * 注意：后端返回的数据通常是倒序（新→旧），需要先反转数组
 */
export function calculateDiskIOSpeed(records: DiskIORecord[]): DiskIOSpeedRecord[] {
  if (records.length < 2) return []

  const result: DiskIOSpeedRecord[] = []

  for (let i = 1; i < records.length; i++) {
    const prev = records[i - 1]
    const curr = records[i]

    // 计算时间差（秒）
    const timeDiff = (new Date(curr.collect_time).getTime() - new Date(prev.collect_time).getTime()) / 1000

    if (timeDiff <= 0) continue

    // 计算速率（bytes/s），确保不为负数
    result.push({
      collect_time: curr.collect_time,
      read_speed: Math.max(0, (curr.read_bytes - prev.read_bytes) / timeDiff),
      write_speed: Math.max(0, (curr.write_bytes - prev.write_bytes) / timeDiff),
    })
  }

  return result
}

/**
 * 将网络IO累计值转换为速率
 * @param records 按时间正序排列的历史记录（旧→新）
 * @returns 速率记录数组
 *
 * 注意：后端返回的数据通常是倒序（新→旧），需要先反转数组
 */
export function calculateNetworkIOSpeed(records: NetworkIORecord[]): NetworkIOSpeedRecord[] {
  if (records.length < 2) return []

  const result: NetworkIOSpeedRecord[] = []

  for (let i = 1; i < records.length; i++) {
    const prev = records[i - 1]
    const curr = records[i]

    // 计算时间差（秒）
    const timeDiff = (new Date(curr.collect_time).getTime() - new Date(prev.collect_time).getTime()) / 1000

    if (timeDiff <= 0) continue

    // 计算速率（bytes/s），确保不为负数
    result.push({
      collect_time: curr.collect_time,
      upload_speed: Math.max(0, (curr.bytes_sent - prev.bytes_sent) / timeDiff),
      download_speed: Math.max(0, (curr.bytes_recv - prev.bytes_recv) / timeDiff),
    })
  }

  return result
}

/**
 * 按设备分组并计算速率
 * @param records 原始IO记录
 * @param deviceKey 设备标识字段名
 * @param calculateFn 速率计算函数
 */
export function groupAndCalculateSpeed<T extends { disk_name?: string; interface_name?: string; collect_time: string }>(
  records: T[],
  deviceKey: 'disk_name' | 'interface_name',
  calculateFn: (records: T[]) => any[]
): Map<string, any[]> {
  const grouped = new Map<string, T[]>()

  // 按设备分组
  records.forEach(record => {
    const device = (record[deviceKey] || 'unknown') as string
    if (!grouped.has(device)) {
      grouped.set(device, [])
    }
    grouped.get(device)!.push(record)
  })

  // 对每个设备分别计算速率
  const result = new Map<string, any[]>()
  grouped.forEach((deviceRecords, device) => {
    // 按时间排序（旧→新）
    const sorted = [...deviceRecords].sort((a, b) =>
      new Date(a.collect_time).getTime() - new Date(b.collect_time).getTime()
    )
    result.set(device, calculateFn(sorted))
  })

  return result
}
