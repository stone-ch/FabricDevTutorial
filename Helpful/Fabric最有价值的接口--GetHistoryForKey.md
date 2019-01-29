## Fabric 最有价值的接口-- GetHistoryForKey

#### 写在前边

文中暂时没有给出具体的实验过程，有时间再补，有经验的读者可以自行实验。我的实验环境 StateDB 为 CouchDB。

#### Fabric 的查询

先说一说 Fabric 的查询，说到查询得先明白一个东西-- StateDB，Fabric 中有一个 StateDB 来维护 “World State” 我们使用的查询接口，无论是 GetState(),还是支持 CouchDB 富查询的接口 GetQueryResult(),都是从 StateDB 中查询的，并非从 ledger 中查询。这里所说的 StateDB 是指 golevelDB 和 CouchDB，ledger 是账本也就是区块文件。

GetHistoryForKey() 的功能是根据 key 查询历史记录。只有这个查询接口才是根据账本的记录从账本中把历史记录都找出来。

#### 为什么 GetHistoryForKey 最有价值

上边说了 GetHistoryForKey 是从账本中查找历史记录，我们知道区块链不可篡改的是账本上的内容，而 World State 是记录在数据库中的，并没有防篡改的机制，所以说只有从账本中查询到的数据才是可信的。所以当需要保证查找到的数据的真实性时最好还是使用 GetHistoryForKey 进行查找。

#### 如何验证

可以在链码中使用 PutState(key, value) 写入一条记录，然后使用 GetState(key) 和 GetHistoryForKey(key) 分别获取记录。然后直接把 CouchDB 中的 key 对应的数据改了，再使用 GetState(key) 和 GetHistoryForKey(key) 查询的时候，你会发现，GetState(key) 获取到的数据是改了 CouchDB 中数据之后的，而 GetHistoryForKey(key) 得到的才是真正的数据。

实验步骤可以参考 (https://www.cnblogs.com/studyzy/p/7101136.html)[https://www.cnblogs.com/studyzy/p/7101136.html] 这篇博文。


