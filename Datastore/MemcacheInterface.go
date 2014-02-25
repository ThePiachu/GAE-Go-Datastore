package Datastore

// Copyright 2012-2014 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import(
	"appengine"
	"appengine/memcache"
	"encoding/gob"
	"appengine/capability"
	"bytes"
	"appengine/datastore"
	"errors"
)

func PutInMemcache(c appengine.Context, key string, toStore interface{}){
	if !capability.Enabled(c, "memcache", "*") {
		c.Errorf("Datastore - PutInMemcache - error 1 - Memcache not available.")
		return
	}
	var data bytes.Buffer
	
	enc:=gob.NewEncoder(&data)
	
	err:=enc.Encode(toStore)
	if err!=nil{
		c.Errorf("Datastore - PutInMemcache - error 2 - %s", err)
		return
	}
	item:=&memcache.Item{
		Key:	key,
		Value:	data.Bytes(),
	}
	if err := memcache.Set(c, item); err != nil {
		c.Errorf("Datastore - PutInMemcache - error 3 - %s", err)
	}
}

func GetFromMemcache(c appengine.Context, key string, dst interface{}) interface{}{
	if !capability.Enabled(c, "memcache", "*") {
		c.Errorf("Datastore - GetFromMemcache - error 1 - Memcache not available.")
		return nil
	}
	item, err := memcache.Get(c, key)
	if err != nil && err != memcache.ErrCacheMiss {
		c.Errorf("Datastore - GetFromMemcache - error 2 - %s", err)
		return nil
	} else if err==memcache.ErrCacheMiss {
		return nil
	}
	dec := gob.NewDecoder(bytes.NewBuffer(item.Value))
	err = dec.Decode(dst)
	if err != nil{
		c.Errorf("Datastore - GetFromMemcache - error 3 - %s", err)
		return nil
	}
	return dst
}

func PutInDatastoreSimpleAndMemcache(c appengine.Context, kind, stringID, memcacheID string, variable interface{}) (*datastore.Key, error){
	if !capability.Enabled(c, "datastore_v3", "*") {
		c.Errorf("Datastore - PutInDatastoreSimpleAndMemcache - error 1 - Datastore not available.")
		return nil, errors.New("Datastore not available")
	}
	
	key, err:=PutInDatastoreSimple(c, kind, stringID, variable)
	if err!=nil{
		c.Errorf("Datastore - PutInDatastoreSimpleAndMemcache - error 2 - %s", err)
		return nil, err
	}
	
	if capability.Enabled(c, "memcache", "*") {
		PutInMemcache(c, memcacheID, variable)
	}
	
	return key, nil
}

func GetFromDatastoreSimpleOrMemcache(c appengine.Context, kind, stringID, memcacheID string, dst interface{}) error{
	if capability.Enabled(c, "memcache", "*") {
		answer:=GetFromMemcache(c, memcacheID, dst)
		if answer!=nil{
			dst=answer
			return nil
		}
	}
	if !capability.Enabled(c, "datastore_v3", "*") {
		c.Errorf("Datastore - GetFromDatastoreOrMemcache - error 1 - Datastore not available.")
		return errors.New("Datastore not available")
	}
	
	err:=GetFromDatastoreSimple(c, kind, stringID, dst)
	if err!=nil{
		c.Infof("Trying to reach - %s - %s", kind, stringID)
		c.Infof("Datastore - GetFromDatastoreOrMemcache error 1 - %s", err)
		return err
	}
	
	if capability.Enabled(c, "memcache", "*") {
		PutInMemcache(c, memcacheID, dst)
	}
	
	return nil
}

func IsVariableInDatastoreSimpleOrMemcache(c appengine.Context, kind, stringID, memcacheID string, dst interface{}) bool{
	_, err := memcache.Get(c, memcacheID)
	if err == nil {
		return true
	}
	return IsVariableInDatastoreSimple(c, kind, stringID, dst)
}

func DeleteFromMemcache(c appengine.Context, memcacheID string){
	memcache.Delete(c, memcacheID)
}

func DeleteFromDatastoreSimpleAndMemcache(c appengine.Context, kind, stringID, memcacheID string)error{
	DeleteFromMemcache(c, memcacheID)
	return DeleteFromDatastoreSimple(c, kind, stringID)
}

func TestMemcache(c appengine.Context){
	type TMP struct{
		A string
		B int
		C float64
	}
	tmp := new(TMP)
	tmp.A="Hello"
	tmp.B=123
	tmp.C=12.3
	c.Infof("TestMemcache")
	key, err:=PutInDatastoreSimpleAndMemcache(c, "test", "test", "test", tmp)
	c.Infof("%v, %v", key, err)
	c.Infof("%v", GetFromDatastoreSimpleOrMemcache(c, "test", "test", "test", tmp))
	c.Infof("%v", tmp)
}