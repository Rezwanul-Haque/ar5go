package serializers

type HealthResp struct {
	DBOnline    bool `json:"db_online"`
	CacheOnline bool `json:"cache_online"`
}
