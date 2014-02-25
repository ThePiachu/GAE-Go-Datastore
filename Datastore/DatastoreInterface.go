package Datastore

// Copyright 2012-2014 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import(
	"appengine"
	"appengine/datastore"
)

func PutInDatastoreFull(c appengine.Context, kind, stringID string, intID int64, parent *datastore.Key, variable interface{}) (*datastore.Key, error){
	k := datastore.NewKey(c, kind, stringID, intID, parent)
	key, err:=datastore.Put(c, k, variable)
	return key, err
}

func PutInDatastoreSimple(c appengine.Context, kind, stringID string, variable interface{})(*datastore.Key, error){
	return PutInDatastoreFull(c, kind, stringID, 0, nil, variable)
}
func PutInDatastore(c appengine.Context, kind string, variable interface{})(*datastore.Key, error){
	return PutInDatastoreFull(c, kind, "", 0, nil, variable)
}

func GetFromDatastoreFull(c appengine.Context, kind, stringID string, intID int64, parent *datastore.Key, dst interface{}) error{
	k := datastore.NewKey(c, kind, stringID, intID, parent)
	return datastore.Get(c, k, dst)
}

func GetFromDatastoreSimple(c appengine.Context, kind, stringID string, dst interface{}) error{
	return GetFromDatastoreFull(c, kind, stringID, 0, nil, dst)
}

//A function that either loads a variable from datastore, or if it is not present, sets it and then loads it
func GetFromDatastoreOrSetDefaultFull(c appengine.Context, kind, stringID string, intID int64, parent *datastore.Key, dst interface{}, def interface{}) error{
	
	key := datastore.NewKey(c, kind, stringID, intID, parent)
	if err := datastore.Get(c, key, dst); err != nil {
		if err.Error()=="datastore: no such entity"{
			_, err2:=datastore.Put(c, key, def)
			
			if err2!=nil{
				return err2
			} else {
				if err3 := datastore.Get(c, key, dst); err3 != nil {
					return err3
				}
			}
			
		} else {
			return err
		}
	}
	return nil
}

func GetFromDatastoreOrSetDefaultSimple(c appengine.Context, kind, stringID string, dst interface{}, def interface{}) error{
	return GetFromDatastoreOrSetDefaultFull(c, kind, stringID, 0, nil, dst, def)
}





func IsVariableInDatastoreSimple(c appengine.Context, kind, stringID string, dst interface{}) bool{
	//var dst *interface{}
	err:=GetFromDatastoreSimple(c, kind, stringID, dst)
	if err==nil{
		return true
	}
	if err.Error()=="datastore: no such entity"{
		return false
	}
	c.Errorf("IsVariableInDatastoreSimple - %s", err)
	return false
}

func QueryGetAllWithFiler(c appengine.Context, kind string, filterStr string, filterValue interface{}, dst interface{})([]*datastore.Key, error){
	return QueryGetAllWithFilerAndLimit(c, kind, filterStr, filterValue, -1, dst)
}

func QueryGetAllWithFilerAndLimit(c appengine.Context, kind string, filterStr string, filterValue interface{}, limit int, dst interface{})([]*datastore.Key, error){
	q:=datastore.NewQuery(kind).Filter(filterStr, filterValue).Limit(limit)
	return q.GetAll(c, dst)
}

func QueryGetAll(c appengine.Context, kind string, dst interface{})([]*datastore.Key, error){
	q:=datastore.NewQuery(kind)
	//q.KeysOnly()
	return q.GetAll(c, dst)
}

func QueryGetAllKeysWithFiler(c appengine.Context, kind string, filterStr string, filterValue interface{}, dst interface{})([]*datastore.Key, error){
	return QueryGetAllKeysWithFilerAndLimit(c, kind, filterStr, filterValue, -1, dst)
}

func QueryGetAllKeysWithFilerAndLimit(c appengine.Context, kind string, filterStr string, filterValue interface{}, limit int, dst interface{})([]*datastore.Key, error){
	q:=datastore.NewQuery(kind).Filter(filterStr, filterValue).Limit(limit).KeysOnly()
	return q.GetAll(c, dst)
}

func QueryGetAllKeys(c appengine.Context, kind string, dst interface{})([]*datastore.Key, error){
	q:=datastore.NewQuery(kind).KeysOnly()
	//q.KeysOnly()
	return q.GetAll(c, dst)
}

func CountQueryWithFilter(c appengine.Context, kind string, filterStr string, filterValue interface{}) int{
	q:=datastore.NewQuery(kind).Filter(filterStr, filterValue)
	count, err:=q.Count(c)
	if err!=nil{
		c.Errorf("CountQueryWithFilter - %s", err)
		return -1
	}
	return count
}

func ClearNamespace(c appengine.Context, kind string){
	q:=datastore.NewQuery(kind)
	q=q.KeysOnly()
	
	keys, err:=q.GetAll(c, nil)
	
	if err!=nil{
		c.Errorf("Clear Namespace 1 - %v", err)
		return
	}
	
	err=datastore.DeleteMulti(c, keys)
	
	if err!=nil{
		c.Errorf("ClearNamespace 2 - %v", err)
	}
}


func DeleteFromDatastoreFull(c appengine.Context, kind, stringID string, intID int64, parent *datastore.Key) error{
	k := datastore.NewKey(c, kind, stringID, intID, parent)
	return datastore.Delete(c, k)
}

func DeleteFromDatastoreSimple(c appengine.Context, kind, stringID string) error{
	return DeleteFromDatastoreFull(c, kind, stringID, 0, nil)
}