并发数：10

超时时间5s ， 握手超时5s 使用pool  2m1s  \pm 5s<br> 命中率约5/7 
超时时间5s ， 握手超时5s  不使用pool  3m38s  \pm 5s

超时时间10s  ， 握手超时10s 使用pool  2m1s   \pm 5s 命中率约6/7 <br>
超时时间10s  ， 握手超时10s 使用pool  3m46   \pm 5s 

并发数：100

超时时间5s ， 握手超时5s 使用pool  50s \pm 5s 命中率约5/7<br>
超时时间5s ， 握手超时5s  不使用pool  100s  \pm 5s

并发数：1000

超时时间5s ， 握手超时5s 使用pool  15s \pm 2s  15% 这位更是拉大的 :clown_face:<br>
超时时间5s ， 握手超时5s  不使用pool  15s \pm 2s

极高并发的情况理论上不太好，同域名可能会有429,超时时间也影响访问的命中率，不过疑似校园网才是更大的客观因素<br>
高并发导致我的网络占用高而访问不到网址（乐 <br>
使用读写分离可能还要更快 <br>
sync.pool 魅力时刻 :tada: