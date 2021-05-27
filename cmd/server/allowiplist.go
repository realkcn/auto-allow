package main

import (
	"github.com/ReneKroon/ttlcache/v2"
	log "github.com/sirupsen/logrus"
	"time"
)

type TTLList struct {
	listCache      *ttlcache.Cache
	removeCallback ttlcache.ExpireReasonCallback
}

func NewTTLList(ttl time.Duration, removeCallback ttlcache.ExpireReasonCallback) *TTLList {
	ttlList := new(TTLList)
	ttlList.listCache = ttlcache.NewCache()
	ttlList.SetTTL(ttl)
	ttlList.removeCallback = removeCallback
	ttlList.listCache.SetExpirationReasonCallback(removeCallback)
	ttlList.listCache.SkipTTLExtensionOnHit(true)
	return ttlList
}

func (ttlList *TTLList) SetTTL(duration time.Duration) {
	err := ttlList.listCache.SetTTL(duration)
	if err != nil {
		log.Errorf("Set TTL error: %s", err.Error())
		return
	}
}

func (ttlList *TTLList) Exist(item string) bool {
	if _, err := ttlList.listCache.Get(item); err == ttlcache.ErrNotFound {
		return false
	}
	return true
}

func (ttlList *TTLList) Add(item string) {
	err := ttlList.listCache.Set(item, "")
	if err != nil {
		log.Errorf("Set %s error: %s", item, err.Error())
		return
	}
}

func (ttlList *TTLList) Remove(item string) {
	err := ttlList.listCache.Remove(item)
	if err != nil {
		log.Errorf("Remove %s error: %s", item, err.Error())
		return
	}
}

func (ttlList *TTLList) GetAll() []string {
	return ttlList.listCache.GetKeys()
}
