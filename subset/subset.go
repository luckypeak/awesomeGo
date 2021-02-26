package subset

import "math/rand"

/*
https://mp.weixin.qq.com/s?__biz=MjM5MDUwNTQwMQ==&mid=2257484430&idx=1&sn=69ead747d0b5cfdcb0cc288790d5c037&chksm=a5391a58924e934ee8f65a79ffee1ccad1ce5b34670cfca3066b9bd06dd244f5a520de209889&scene=132#wechat_redirect
  构造本地连接池算法，在微服务场景下，单一服务可能有上万个实例，每个client 和服务都建立连接，将会占用大量连接。
   需要构造连接池，同时保证1. 每个服务连接数是均衡的，2。上下线不能有过多迁移，3. 连接数要足够
   使用谷歌的subset 算法，有以下限制 1. clientID 要自增有序，2. 一次不能上下线过多机器

 */
func Subset(backends []string, clientID int, subsetSize int) []string{
	subSetCount := len(backends) /subsetSize
	round := clientID/subSetCount
	r := rand.New(rand.NewSource(int64(round)))
	r.Shuffle(len(backends), func(i, j int) {
		backends[i], backends[j] = backends[j], backends[i]
	})
	subSetId := clientID % subSetCount
	start := subSetId * subsetSize
	return backends[start: start+subsetSize]
}
