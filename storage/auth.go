package storage

/*
User is used in user system.
*/
// type User struct {
// 	ID  string `yaml:"id" json:"id"`
// 	Key string `yaml:"key" json:"key"`
// }

// var authUser *User
// var userMap map[string]User

// func InitUserDb() error {
// 	userMap = make(map[string]User)
// 	cfg := config.GetConfig()
// 	for _, user := range cfg.Users {
// 		userMap[user.Key] = User{
// 			ID:  user.ID,
// 			Key: user.Key,
// 		}
// 		log.Printf("- user %d: %s\n", user.ID, user.Key)
// 	}
// 	return nil
// }

// func getUser(userId string, key string) *User {
// 	if val, ok := userMap[key]; ok {
// 		if val.ID == userId {
// 			return &val
// 		}
// 	}
// 	return nil
// }

// func AuthUser(userId string, key string) *User {
// 	if user := getUser(userId, key); user != nil {
// 		return user
// 	}
// 	return nil
// }
