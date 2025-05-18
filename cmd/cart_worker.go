package main

// import (
//     "context"
//     "log"
//     "strings"
//     "time"

//     "rijig/config"
//     "rijig/internal/services"
// )

// // func main() {
// //     config.SetupConfig()

// // }

// func processCartKeys(ctx context.Context, cartService services.CartService) {
//     pattern := "cart:user:*"
//     iter := config.RedisClient.Scan(ctx, 0, pattern, 0).Iterator()

//     for iter.Next(ctx) {
//         key := iter.Val()
//         ttl, err := config.RedisClient.TTL(ctx, key).Result()
//         if err != nil {
//             log.Printf("Failed to get TTL for key %s: %v", key, err)
//             continue
//         }

//         if ttl <= time.Minute {
//             log.Printf("🔄 Auto-committing key: %s", key)
//             parts := strings.Split(key, ":")
//             if len(parts) != 3 {
//                 log.Printf("Invalid key format: %s", key)
//                 continue
//             }
//             userID := parts[2]

//             err := cartService.CommitCartFromRedis(userID)
//             if err != nil {
//                 log.Printf("❌ Failed to commit cart for user %s: %v", userID, err)
//             } else {
//                 log.Printf("✅ Cart for user %s committed successfully", userID)
//             }
//         }
//     }

//     if err := iter.Err(); err != nil {
//         log.Printf("Error iterating keys: %v", err)
//     }
// }
