package service

import (
  "fmt"
  "os"
  "testing"
)

func TestOne(t *testing.T) {
  os.Setenv("DB_URL", "cailianpress_dba:xxxxx@tcp(localhost:3306)/opt?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai")
  os.Setenv("REDIS_ADDR",  "localhost:6379")
  os.Setenv("REDIS_PWD",  "xxxxx")

  OpenDB()
  defer CloseDB()

  OpenRedis()
  defer CloseRedis()

  const key = "mq:c1"
  redisCache.LPush(key, "1")

  fmt.Println("1---", redisCache.RPop(key).Val())
  fmt.Println("2---",redisCache.RPop(key).Val())

}

func change(arr *[]int, i *int, m *map[int]string) (*[]int, *int, *map[int]string) {
  // arr = &[]int{1,2,3}
  // arr[0] = 1
  *arr = append(*arr, 1)
  *arr = append(*arr, 2)
  //k := 2
  *i = 2
  fmt.Printf("k = %p \n", i)
  return arr, i, m
}

func quickSort(values []int, left, right int) {
  temp := values[left]
  p := left
  i, j := left, right
  for i <= j {
    for j >= p && values[j] >= temp {
      j--
    }
    if j >= p {
      values[p] = values[j]
      p = j
    }
    for i <= p && values[i] <= temp {
      i++
    }
    if i <= p {
      values[p] = values[i]
      p = i
    }
  }

  values[p] = temp

  if p-left > 1 {
    quickSort(values, left, p-1)
  }
  if right-p > 1 {
    quickSort(values, p+1, right)
  }

}

func QuickSort(values []int) {
  quickSort(values, 0, len(values)-1)
}
