package group

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisGroupStore Redis 分组存储
type RedisGroupStore struct {
	client redis.UniversalClient
	prefix string
	ttl    time.Duration
}

// NewRedisGroupStore 创建 Redis 分组存储
func NewRedisGroupStore(client redis.UniversalClient, prefix string, ttl time.Duration) *RedisGroupStore {
	return &RedisGroupStore{
		client: client,
		prefix: prefix,
		ttl:    ttl,
	}
}

// groupKey 生成分组键
func (r *RedisGroupStore) groupKey(groupID string) string {
	return fmt.Sprintf("%s:group:%s", r.prefix, groupID)
}

// groupMetaKey 生成分组元数据键
func (r *RedisGroupStore) groupMetaKey(groupID string) string {
	return fmt.Sprintf("%s:group_meta:%s", r.prefix, groupID)
}

// clientGroupsKey 生成客户端分组键
func (r *RedisGroupStore) clientGroupsKey(clientID string) string {
	return fmt.Sprintf("%s:client_groups:%s", r.prefix, clientID)
}

// AddClientToGroup 添加客户端到分组
func (r *RedisGroupStore) AddClientToGroup(ctx context.Context, groupID, clientID string, clientInfo map[string]interface{}) error {
	pipe := r.client.Pipeline()

	// 将客户端添加到分组
	clientInfoBytes, err := json.Marshal(clientInfo)
	if err != nil {
		return err
	}
	pipe.HSet(ctx, r.groupKey(groupID), clientID, clientInfoBytes)
	pipe.Expire(ctx, r.groupKey(groupID), r.ttl)

	// 将分组添加到客户端的分组列表
	pipe.SAdd(ctx, r.clientGroupsKey(clientID), groupID)
	pipe.Expire(ctx, r.clientGroupsKey(clientID), r.ttl)

	_, err = pipe.Exec(ctx)
	return err
}

// RemoveClientFromGroup 从分组移除客户端
func (r *RedisGroupStore) RemoveClientFromGroup(ctx context.Context, groupID, clientID string) error {
	pipe := r.client.Pipeline()

	// 从分组中移除客户端
	pipe.HDel(ctx, r.groupKey(groupID), clientID)

	// 从客户端的分组列表中移除分组
	pipe.SRem(ctx, r.clientGroupsKey(clientID), groupID)

	// 如果分组为空，删除分组元数据
	pipe.HLen(ctx, r.groupKey(groupID))

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	// 检查分组是否为空
	if len(cmds) >= 3 {
		if hlenCmd, ok := cmds[2].(*redis.IntCmd); ok {
			if hlenCmd.Val() == 0 {
				// 分组为空，删除元数据
				r.client.Del(ctx, r.groupMetaKey(groupID))
			}
		}
	}

	return nil
}

// GetGroupClients 获取分组中的所有客户端
func (r *RedisGroupStore) GetGroupClients(ctx context.Context, groupID string) (map[string]map[string]interface{}, error) {
	// 获取分组中的所有客户端
	clientsMap, err := r.client.HGetAll(ctx, r.groupKey(groupID)).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string]map[string]interface{})
	for clientID, clientInfoBytes := range clientsMap {
		var clientInfo map[string]interface{}
		if err := json.Unmarshal([]byte(clientInfoBytes), &clientInfo); err != nil {
			continue
		}
		result[clientID] = clientInfo
	}

	return result, nil
}

// GetClientGroups 获取客户端加入的所有分组
func (r *RedisGroupStore) GetClientGroups(ctx context.Context, clientID string) ([]string, error) {
	return r.client.SMembers(ctx, r.clientGroupsKey(clientID)).Result()
}

// SetGroupMetadata 设置分组元数据
func (r *RedisGroupStore) SetGroupMetadata(ctx context.Context, groupID string, metadata map[string]interface{}) error {
	// 将元数据序列化为 JSON
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	// 设置分组元数据
	return r.client.Set(ctx, r.groupMetaKey(groupID), metadataBytes, r.ttl).Err()
}

// GetGroupMetadata 获取分组元数据
func (r *RedisGroupStore) GetGroupMetadata(ctx context.Context, groupID string) (map[string]interface{}, error) {
	// 获取分组元数据
	metadataBytes, err := r.client.Get(ctx, r.groupMetaKey(groupID)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var metadata map[string]interface{}
	if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

// GetAllGroups 获取所有分组
func (r *RedisGroupStore) GetAllGroups(ctx context.Context) ([]string, error) {
	// 使用 Scan 命令查找所有分组键
	var groups []string
	var cursor uint64
	pattern := fmt.Sprintf("%s:group:*", r.prefix)

	for {
		var keys []string
		var err error
		keys, cursor, err = r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			// 从键中提取分组 ID
			groupID := key[len(r.prefix)+7:] // 去掉前缀 "{prefix}:group:"
			groups = append(groups, groupID)
		}

		if cursor == 0 {
			break
		}
	}

	return groups, nil
}

// GetGroupSize 获取分组大小
func (r *RedisGroupStore) GetGroupSize(ctx context.Context, groupID string) (int64, error) {
	return r.client.HLen(ctx, r.groupKey(groupID)).Result()
}

// GroupExists 检查分组是否存在
func (r *RedisGroupStore) GroupExists(ctx context.Context, groupID string) (bool, error) {
	exists, err := r.client.Exists(ctx, r.groupKey(groupID)).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// ClientInGroup 检查客户端是否在分组中
func (r *RedisGroupStore) ClientInGroup(ctx context.Context, groupID, clientID string) (bool, error) {
	exists, err := r.client.HExists(ctx, r.groupKey(groupID), clientID).Result()
	if err != nil {
		return false, err
	}
	return exists, nil
}

// DeleteGroup 删除分组
func (r *RedisGroupStore) DeleteGroup(ctx context.Context, groupID string) error {
	// 获取分组中的所有客户端
	clientIDs, err := r.client.HKeys(ctx, r.groupKey(groupID)).Result()
	if err != nil {
		return err
	}

	pipe := r.client.Pipeline()

	// 从每个客户端的分组列表中移除分组
	for _, clientID := range clientIDs {
		pipe.SRem(ctx, r.clientGroupsKey(clientID), groupID)
	}

	// 删除分组和元数据
	pipe.Del(ctx, r.groupKey(groupID))
	pipe.Del(ctx, r.groupMetaKey(groupID))

	_, err = pipe.Exec(ctx)
	return err
}
