package rdb

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Task struct {
	delayedTasksKey    string
	processingTasksKey string
	retryTasksKey      string
	maxRetryCount      int64
}

func NewTask(taskKey string, maxRetryCount int64) *Task {
	return &Task{
		delayedTasksKey:    "delayed:" + taskKey,
		processingTasksKey: "processing:" + taskKey,
		retryTasksKey:      "retry:" + taskKey,
		maxRetryCount:      maxRetryCount,
	}
}

func (t *Task) Add(ctx context.Context, prefix, key string, executeAt int64) error {
	now := time.Now().Unix()
	if executeAt <= now {
		executeAt += now
	}
	return rdb.ZAdd(ctx, t.delayedTasksKey, redis.Z{Score: float64(executeAt), Member: prefix + ":" + key}).Err()
}

func (t *Task) Process(ctx context.Context, f func(string) error) {
	for {
		now := time.Now().Unix()
		tasks, err := rdb.ZRangeByScoreWithScores(ctx, t.delayedTasksKey,
			&redis.ZRangeBy{Min: "0", Max: strconv.FormatInt(now, 10), Offset: 0, Count: 1}).
			Result()
		if err != nil {
			zap.L().Error("Failed to fetch tasks", zap.Error(err))
			continue
		}
		for _, task := range tasks {
			key := task.Member.(string)
			// 移除任务
			rdb.ZRem(ctx, t.delayedTasksKey, key)
			// 执行任务
			zap.L().Info("Executing task", zap.String("key", key), zap.Int64("score", now))
			if err := f(key); err != nil {
				zap.L().Error("Failed to execute task", zap.String("key", key), zap.Float64("score", task.Score), zap.Error(err))
				// 重试
				count, err := rdb.Incr(ctx, t.retryTasksKey+key).Result()
				if err != nil {
					zap.L().Error("Failed to increase retry count", zap.String("key", key), zap.Error(err))
				}
				if count <= t.maxRetryCount {
					// 重试延迟 count*10 秒
					rdb.ZAdd(ctx, t.delayedTasksKey, redis.Z{Score: float64(now + 10*count), Member: key})
				} else {
					zap.L().Error("Retry count exceeded", zap.String("key", key), zap.Int64("count", count))
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func (t *Task) Start(ctx context.Context, f func(string) error) {
	go t.Process(ctx, f)
}

func (t *Task) Remove(ctx context.Context, key string) error {
	members, err := rdb.ZRange(ctx, t.delayedTasksKey, 0, -1).Result()
	if err != nil {
		return err
	}
	delMembers := make([]string, 0)
	for _, member := range members {
		if strings.Contains(member, key) {
			delMembers = append(delMembers, member)
		}
	}
	return rdb.ZRem(ctx, t.delayedTasksKey, delMembers).Err()
}

type TaskItems []*TaskItem

type TaskItem struct {
	Key   string
	Score int64
}

func (t *TaskItems) String() string {
	buf := bytes.Buffer{}
	for i, item := range *t {
		buf.WriteString(fmt.Sprintf("%d. %s %s\n", i+1, item.Key, time.Unix(item.Score, 0).Format("2006-01-02 15:04:05")))
	}
	return buf.String()
}

func (t *Task) List(ctx context.Context) (TaskItems, error) {
	r, err := rdb.ZRangeWithScores(ctx, t.delayedTasksKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	items := make([]*TaskItem, 0, len(r))
	for _, v := range r {
		items = append(items, &TaskItem{Key: v.Member.(string), Score: int64(v.Score)})
	}
	return items, nil
}
