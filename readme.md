# repo-image-hosting 🐽

> Github / Gitee 图床工具 (基于 Golang(Gin) 实现) [github地址](https://github.com/hezhizheng/repo-image-hosting)

```
还是用Gitee当图床算了，不然哪里有国内访问又快又免费又稳的图床服务提供......
Github 也行，可以用jsdelivr加速
```

## 页面
![](https://gitee.com/hezhizheng/pictest/raw/master/newdir/20210207160931_ZPUDMEHPYYLWOQFJ.png)
![](https://cdn.learnku.com/uploads/images/202102/07/6843/8CY2HIkX5x.gif!large)

## 功能
- 支持使用 Github / Gitee 作为图床工具
- github 使用 jsdelivr 加速 
- 支持gitee图片大于1M的显示(PS：需要手动启动/部署gitee的pages服务)   
- 一键启动，跨平台支持，运行只依赖编译后的二进制文件
- 可视化web操作界面(PS: 页面有点丑，不影响操作......)
- 多图上传，不限制图片类型
- 复制图片url 、删除图片

## 使用
用户可直接下载 [releases](https://github.com/hezhizheng/repo-image-hosting/releases) 文件启动即可，参数说明：
![](https://cdn.jsdelivr.net/gh/hezhizheng/static-image-hosting@master/image-hosting/20210222093638_VUQZUKNZGAXXSXJI.png)

```
./repo-image-hosting_windows_amd64.exe -h
Usage of D:\phpstudy_pro\WWW\org\gitee-image-hosting\repo-image-hosting_windows_amd64.exe:
  -owner string
        仓库所属空间地址(企业、组织或个人的地址path) (default "hezhizheng")
  -path string
        文件的路径 (default "image-hosting")
  -platform string
        平台名称，支持gitee/github (default "github")
  -port string
        本地监听的端口 (default "2047")
  -repo string
        仓库路径(path) (default "static-image-hosting")
  -token string
        Gitee/Github 的用户授权码
```

```
完整启动命令： ./repo-image-hosting_windows_amd64.exe -platform github -owner hezhizheng -repo static-image-hosting -path image-hosting -token xxxtoken -port 2047
实际参数替换成自己的就行
```

token获取(gitee)：https://gitee.com/profile/personal_access_tokens/new

token获取(github)：https://github.com/settings/tokens/new

自行编译
```
// 编译前请先安装 go-bindata ，请参考 https://blog.hi917.com/detail/87.html
// 执行静态资源编译命令（每次修改静态文件都需要执行）
go-bindata -o=bindata/bindata.go -pkg=bindata ./static/... ./views/... 

// 跨平台编译
gox -osarch="windows/amd64" -ldflags "-s -w" -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}"

gox -osarch="darwin/amd64" -ldflags "-s -w" -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}"

gox -osarch="linux/amd64" -ldflags "-s -w" -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}"
```


## 关于Gitee限制图片大于1M访问的处理方案
- 使用第三方图片压缩工具进行压缩，之后再进行上传。推荐 [compressjpeg](https://compressjpeg.com/zh/)
- 启用Gitee的pages功能(非付费用户上传图片之后需要手动进行pages服务的部署)，程序会自动替换pages域名进行图片的展示。


## License
[MIT](./LICENSE.txt)