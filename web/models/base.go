package models

import (
	"yulong-hids/web/models/wmongo"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

// baseModel 基础模型
type baseModel struct {
	collectionName string
}

// GetAll get all document in collection
func (bmodel *baseModel) GetAll() []bson.M {
	return bmodel.GetPieces(nil, 0, 0)
}

//GetPieces limit and start in find
func (bmodel *baseModel) GetPieces(q bson.M, start int, limit int) []bson.M {
	return bmodel.Find(q, start, limit)
}

// GetSortedTop find for sort
func (bmodel *baseModel) GetSortedTop(q bson.M, start int, limit int, sort ...string) []bson.M {
	return bmodel.Find(q, start, limit, sort...)
}

// FindOne get one result
func (bmodel *baseModel) FindOne(query bson.M) bson.M {
	res := bmodel.Find(query, 0, 1)
	if len(res) == 1 {
		return res[0]
	}
	return nil
}

// FindByID as name
func (bmodel *baseModel) FindByID(id bson.ObjectId) bson.M {
	res := bmodel.FindOne(bson.M{"_id": id})
	return res
}

// FindAll as name
func (bmodel *baseModel) FindAll(query bson.M) []bson.M {
	res := bmodel.Find(query, 0, 0)
	return res
}

// Find base find function
func (bmodel *baseModel) Find(query bson.M, start int, limit int, sort ...string) []bson.M {

	mConn := wmongo.Conn()
	defer mConn.Close()
	cname := bmodel.collectionName
	collections := mConn.DB("").C(cname)

	var cli []bson.M

	if len(sort) == 0 {
		sort = append(sort, "-uptime")
	}

	if err := collections.Find(query).Skip(start).Limit(limit).Sort(sort...).All(&cli); err != nil {
		beego.Error("Model Find Error:", err, query, start, limit, sort)
	}

	return cli
}

// Count count the file result number
func (bmodel *baseModel) Count(query bson.M) int {
	mConn := wmongo.Conn()
	defer mConn.Close()
	cname := bmodel.collectionName
	collections := mConn.DB("").C(cname)

	var count int
	var err error
	if count, err = collections.Find(query).Count(); err != nil {
		beego.Error("Model Count(collections.Find) Error", err, query)
	}
	beego.Debug("Show collectionName, query, count:", bmodel.collectionName, query, count)
	return count
}

// CountSubList count the list item in one collection
func (bmodel *baseModel) CountSubList(query bson.M, listkey string) int {

	var count int
	unwind := bson.M{"$unwind": "$" + listkey}
	group := bson.M{"$group": bson.M{"_id": "", "count": bson.M{"$sum": 1}}}
	slice := bmodel.Aggregate(unwind, group)
	if len(slice) > 0 {
		count = slice[0]["count"].(int)
	} else {
		count = 0
	}
	return count

}

// Aggregate db.aggregate
func (bmodel *baseModel) Aggregate(querylist ...bson.M) []bson.M {

	mConn := wmongo.Conn()
	defer mConn.Close()
	cname := bmodel.collectionName
	collections := mConn.DB("").C(cname)

	var res []bson.M
	pipe := collections.Pipe(querylist)
	err := pipe.All(&res)

	if err != nil {
		beego.Error("Collections pipe(pipe.All) aggregate all", err)
	}

	return res
}

// InsertOne as name
func (bmodel *baseModel) InsertOne(data bson.M) bool {
	mConn := wmongo.Conn()
	defer mConn.Close()
	cname := bmodel.collectionName
	collections := mConn.DB("").C(cname)
	err := collections.Insert(data)
	if err != nil {
		beego.Error("Mongodb insert(collections.Insert) error", err)
		return false
	}
	return true
}

// InsertMany insert a list to mongo
func (bmodel *baseModel) InsertMany(datalist []interface{}) error {
	mConn := wmongo.Conn()
	defer mConn.Close()
	cname := bmodel.collectionName
	collections := mConn.DB("").C(cname)
	err := collections.Insert(datalist...)

	return err
}

// Distinct return only one key in find result
func (bmodel *baseModel) Distinct(query bson.M, key string) []interface{} {
	mConn := wmongo.Conn()
	defer mConn.Close()
	cname := bmodel.collectionName
	collections := mConn.DB("").C(cname)

	var res []interface{}

	err := collections.Find(query).Distinct(key, &res)
	if err != nil {
		beego.Error("Collections Distinct(or Find) Error", err)
	}
	return res
}

// UpdateByID as name
func (bmodel *baseModel) UpdateByID(id interface{}, data interface{}) error {
	mConn := wmongo.Conn()
	defer mConn.Close()
	cname := bmodel.collectionName
	collection := mConn.DB("").C(cname)

	err := collection.UpdateId(id, bson.M{"$set": data})

	return err
}

// Remove remove all conform to query, but error will stop
// TODO it is hardly to supported continue_on_error
func (bmodel *baseModel) Remove(query bson.M) error {
	mConn := wmongo.Conn()
	defer mConn.Close()
	cname := bmodel.collectionName
	collection := mConn.DB("").C(cname)

	_, err := collection.RemoveAll(query)

	return err
}
