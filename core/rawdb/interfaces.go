
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:34</date>
//</624342617135517696>


package rawdb

//DatabaseReader包装支持数据存储的has和get方法。
type DatabaseReader interface {
	Has(key []byte) (bool, error)
	Get(key []byte) ([]byte, error)
}

//databasewriter包装支持数据存储的Put方法。
type DatabaseWriter interface {
	Put(key []byte, value []byte) error
}

//databaseDeleter包装备份数据存储的删除方法。
type DatabaseDeleter interface {
	Delete(key []byte) error
}

