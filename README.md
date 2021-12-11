# qqtaotao
QQ说说全量爬取

# 食用方式
1. 登录你的qq空间
2. F12打开，随便触发一次请求，然后把以下的信息复制出来
    ```shell script
    cookie          (header中)
    referer         (header中)
    qzonetoken      (url中)
    g_tk            (url中)
    ```
3. 将上诉的信息以及你的qq号配置在如下代码中
    ```go
    // use to get friends
    var (
        cookie 	= ""
        referer	= ""
        qzonetoken = ""
        g_tk  = ""
        myUin = ""
    )
    ```
4. 访问任意一个好友的说说，把以下信息复制出来
    ```shell script
    cookie          (header中)
    referer         (header中)
    g_tk            (url中)
    ```
5. 将上诉的信息配置在如下代码中
    ```go
    // use to get taotao
    var (
        cookie2 	= ""
        referer2	= ""
        g_tk2  = ""
    )
    ```
6. run

# 食用效果
> 还有很多其他信息可获取，可以自己debug修改

[![o7ArwV.png](https://s4.ax1x.com/2021/12/11/o7ArwV.png)](https://imgtu.com/i/o7ArwV)