## 1. 数据相关

### 1.1 极值数据
- **接口说明：** 极值数据
- **接口方法：** POST
- **接口地址：** /v1/data/extremum
- **参数类型：** JSON

#### 1.1.1 请求头
参数名称						    |参数值		        |描述
:----						    |:---		        |:---
&emsp;Uber-Trace-Id				|Uber-Trace-Id		|open tracing 格式

#### 1.1.1 请求参数

参数名称						    |类型		|出现要求	    |描述
:----						    |:---		|:------	|:---
&emsp;mac				        |[]string	|R			|mac列表
&emsp;metric				    |[]string	|O			|查询的属性 temperature等
&emsp;start_time				|int		|O			|开始时间
&emsp;end_time		            |int		|O			|结束时间

请求示例：

```json5
// https://127.0.0.1:8288/v1/data/extremum
// json body
{
  "mac": ["582D344C0A24","582D344C0A6C"],
  "metric":["temperature","humidity"],
  "start_time": 1620363590,
  "end_time": 1620374391
}
```

#### 1.1.2 返回结果

参数名称						                |类型		|出现要求	|描述
:----						                |:---		|:------	|:---	
code						                |int		|R			|响应码，代码定义请见“[响应吗说明](https://github.com/ClearGrass/QingDaily/blob/master/doc/code.md)”
msg						                    |string		|R			|描述信息
timestamp						            |int		|R			|响应时间
traceId						                |string		|R			|traceId
data						                |object		|R			|数据正文

示例：

```json5
{
    "code": 200,
    "msg": "success",
    "data": {
        "extremum": [
            {
                "mac": "582D344C0424",
                "temperature": {
                    "max": 26.7,
                    "min": 26.1
                },
                "humidity": {
                    "max": 21.3,
                    "min": 20.7
                }
            }
        ]
    }
}
```