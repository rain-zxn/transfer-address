1.编译
```go
go build .
```
2.创建钱包
- count: 创建地址数
- 生成accounts文件夹,其中两个文件,一个为助记词,一个为地址
```
./transfer-address batchcreateaccount --count 30
```

3.给第0个地址转usdt和matic
- 计算一下每个地址12U和1个matic
- 30个地址,需要360U和30.2个matic,分钱时也需要matic,所以多打0.2个

4.分matic和usdt
- 将第0个地址的钱给第[starttosubaccount,--endtosubaccount]的地址
- token: usdt和matic
- tokenamount:数量(必须为整数)

```go
转matic
./transfer-address batchtransfertoken --starttosubaccount 1 --endtosubaccount 29  --token matic --tokenamount 1 --estimate
```

 ```go
转usdt
./transfer-address batchtransfertoken --starttosubaccount 1 --endtosubaccount 29  --token usdt --tokenamount 12 --estimate
```