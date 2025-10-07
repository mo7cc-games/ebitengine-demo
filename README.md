# ebitengine-demo

基于 Ebitengine + Donburi 的 ECS 架构的游戏 Demo。

运行说明:

- 依赖由 go modules 管理，已在 `go.mod` 中声明。
- 在本地运行示例:
  1.  在仓库目录运行 `go mod tidy`（如果尚未执行）以下载依赖。
  2.  运行 `go run .` 或 `go build` 然后执行二进制。

控制:

- 使用方向键或 WASD 移动屏幕中心的蓝色小方块（移动其实是移动视点，使地图看起来是无限的）。
- 按 Esc 退出。

注意：窗口创建与显示依赖于本地图形环境，远程或无头环境可能无法打开窗口。
