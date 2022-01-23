package main

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

func CreateRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   10,
		IdleTimeout: 200 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.DialURL("redis://localhost:6379")
			if err != nil {
				return nil, err
			}
			return conn, err
		},
	}
}

func test10byte(redisPool *redis.Pool) {
	c1 := redisPool.Get()
	var byteTen [10]byte
	for i := 0; i < 10; i++ {
		v := i
		byteTen[i] = byte(v)
	}

	rec0, err := c1.Do("set", "10byte", byteTen)
	fmt.Println(rec0)

	str, err := redis.String(c1.Do("get", "10byte"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(str)
}

func test20byte(redisPool *redis.Pool) {
	c1 := redisPool.Get()
	var byteTwenty [20]byte
	for i := 0; i < 20; i++ {
		v := i
		byteTwenty[i] = byte(v)
	}

	rec0, _ := c1.Do("set", "20byte", byteTwenty)
	fmt.Println(rec0)

	str, err := redis.String(c1.Do("get", "20byte"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(str)
}

func test50byte(redisPool *redis.Pool) {
	c1 := redisPool.Get()
	var byteFifty [50]byte
	for i := 0; i < 50; i++ {
		v := i
		byteFifty[i] = byte(v)
	}

	rec0, _ := c1.Do("set", "50byte", byteFifty)
	fmt.Println(rec0)

	str, err := redis.String(c1.Do("get", "50byte"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(str)
}

func test5000byte(redisPool *redis.Pool) {
	c1 := redisPool.Get()
	var bytess [5000]byte
	for i := 0; i < 5000; i++ {
		v := i
		bytess[i] = byte(v)
	}
	for i := 0; i < 10000; i++ {
		rec0, _ := c1.Do("set", fmt.Sprintf("bytes_no_%d", i), bytess)
		fmt.Println(rec0)
	}
}

//create demo data
func main() {
	redisPool := CreateRedisPool()
	defer redisPool.Close()
	test10byte(redisPool)
	test20byte(redisPool)
	test50byte(redisPool)
	test5000byte(redisPool)

}
