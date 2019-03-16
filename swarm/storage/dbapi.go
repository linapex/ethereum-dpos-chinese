
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:49</date>
//</624342680737943552>

//
//
//
//
//
//
//
//
//
//
//
//
//
//
//

package storage

import "context"

//
type DBAPI struct {
	db  *LDBStore
	loc *LocalStore
}

func NewDBAPI(loc *LocalStore) *DBAPI {
	return &DBAPI{loc.DbStore, loc}
}

//
func (d *DBAPI) Get(ctx context.Context, addr Address) (*Chunk, error) {
	return d.loc.Get(ctx, addr)
}

//
func (d *DBAPI) CurrentBucketStorageIndex(po uint8) uint64 {
	return d.db.CurrentBucketStorageIndex(po)
}

//
func (d *DBAPI) Iterator(from uint64, to uint64, po uint8, f func(Address, uint64) bool) error {
	return d.db.SyncIterator(from, to, po, f)
}

//
func (d *DBAPI) GetOrCreateRequest(ctx context.Context, addr Address) (*Chunk, bool) {
	return d.loc.GetOrCreateRequest(ctx, addr)
}

//
func (d *DBAPI) Put(ctx context.Context, chunk *Chunk) {
	d.loc.Put(ctx, chunk)
}

