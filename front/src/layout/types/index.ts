// Layout 类型定义

export interface NavItem {
  path: string
  label: string
  icon: string
}

export interface UserInfo {
  name: string
  avatar?: string
}

export type LayoutMode = 'default' | 'compact' | 'full'
