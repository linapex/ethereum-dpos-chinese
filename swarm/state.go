
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:49</date>
//</624342680201072640>

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

package swarm

type Voidstore struct {
}

func (self Voidstore) Load(string) ([]byte, error) {
	return nil, nil
}

func (self Voidstore) Save(string, []byte) error {
	return nil
}

