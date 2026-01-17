#!/bin/bash
#
# 安装 Docker 和 Docker Compose 脚本
# 支持 Ubuntu/Debian 和 CentOS/RHEL 系统
#

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检测操作系统类型
detect_os() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$ID
        VERSION=$VERSION_ID
    else
        log_error "无法检测操作系统类型"
        exit 1
    fi
    log_info "检测到操作系统: $OS $VERSION"
}

# 卸载旧版本
uninstall_old() {
    log_info "检查并卸载旧版本 Docker..."
    if command -v docker &> /dev/null; then
        case $OS in
            ubuntu|debian)
                sudo apt-get remove -y docker docker-engine docker.io containerd runc 2>/dev/null || true
                ;;
            centos|rhel|fedora)
                sudo yum remove -y docker docker-client docker-client-latest docker-common \
                    docker-latest docker-latest-logrotate docker-logrotate docker-engine 2>/dev/null || true
                ;;
        esac
        log_info "旧版本 Docker 已卸载"
    else
        log_info "未检测到旧版本 Docker"
    fi
}

# 安装 Docker
install_docker() {
    log_info "开始安装 Docker..."

    case $OS in
        ubuntu|debian)
            # Ubuntu/Debian 安装
            sudo apt-get update
            sudo apt-get install -y \
                ca-certificates \
                curl \
                gnupg \
                lsb-release

            # 添加 Docker 官方 GPG 密钥
            sudo mkdir -p /etc/apt/keyrings
            curl -fsSL https://download.docker.com/linux/$OS/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

            # 设置 Docker 仓库
            echo \
                "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/$OS \
                $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

            # 安装 Docker
            sudo apt-get update
            sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

            # 启动 Docker
            sudo systemctl start docker
            sudo systemctl enable docker

            # 添加当前用户到 docker 组
            sudo usermod -aG docker $USER 2>/dev/null || true
            log_warn "请重新登录或执行 'newgrp docker' 以使用户组生效"

            ;;

        centos|rhel|fedora)
            # CentOS/RHEL 安装
            sudo yum install -y yum-utils
            sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo

            # 安装 Docker
            sudo yum install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

            # 启动 Docker
            sudo systemctl start docker
            sudo systemctl enable docker

            # 添加当前用户到 docker 组
            sudo usermod -aG docker $USER 2>/dev/null || true
            log_warn "请重新登录或执行 'newgrp docker' 以使用户组生效"

            ;;

        *)
            log_error "不支持的操作系统: $OS"
            exit 1
            ;;
    esac

    log_info "Docker 安装完成"
}

# 验证安装
verify_installation() {
    log_info "验证安装..."

    # 检查 Docker
    if command -v docker &> /dev/null; then
        docker_version=$(docker --version)
        log_info "Docker 已安装: $docker_version"

        # 测试 Docker 运行
        sudo docker run --rm hello-world &> /dev/null
        if [ $? -eq 0 ]; then
            log_info "Docker 运行正常"
        else
            log_warn "Docker 测试运行失败，可能需要手动配置"
        fi
    else
        log_error "Docker 安装失败"
        exit 1
    fi

    # 检查 Docker Compose
    if command -v docker-compose &> /dev/null; then
        compose_version=$(docker-compose --version)
        log_info "Docker Compose 已安装: $compose_version"
    elif docker compose version &> /dev/null; then
        compose_plugin_version=$(docker compose version)
        log_info "Docker Compose 插件已安装: $compose_plugin_version"
    else
        log_error "Docker Compose 安装失败"
        exit 1
    fi
}

# 安装完成
install_complete() {
    echo ""
    log_info "=========================================="
    log_info "Docker 和 Docker Compose 安装完成！"
    log_info "=========================================="
    echo ""
    log_info "常用命令:"
    echo "  docker --version           # 查看 Docker 版本"
    echo "  docker ps                  # 查看运行中的容器"
    echo "  docker images              # 查看镜像列表"
    echo "  docker-compose up -d       # 启动服务"
    echo "  docker compose version     # 查看 Docker Compose 版本"
    echo ""
    log_warn "注意: 请重新登录或执行 'newgrp docker' 以使 docker 用户组生效"
}

# 主函数
main() {
    echo ""
    log_info "开始安装 Docker 和 Docker Compose..."
    echo ""

    detect_os
    uninstall_old
    install_docker
    verify_installation
    install_complete

    echo ""
    log_info "安装成功！"
}

# 执行主函数
main
