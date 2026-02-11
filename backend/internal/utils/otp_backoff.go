package utils

import (
	"fmt"
	"time"
	"context"
	"github.com/redis/go-redis/v9"
)


// helper function to implement otp backoff using redis
var rctx=context.Background()


func HandleOtpResendBackoff(
	rdb *redis.Client,
	keyBase string,
)error{

	attemptKey := fmt.Sprintf("%s:attempts", keyBase)
	cooldownKey := fmt.Sprintf("%s:cooldown", keyBase)
	blockKey := fmt.Sprintf("%s:block", keyBase)

	// check if user is blocked
	blocked,err:=rdb.Exists(rctx,blockKey).Result()

	if err!=nil{
		return err
	}

	if blocked==1{
		return fmt.Errorf("too many attempts. try again after 1 hour")
		
	}

	// check cooldown before sending 

	coolDownExist,err:=rdb.Exists(rctx,cooldownKey).Result()

	if err!=nil{
		return err
	}

	if coolDownExist==1{
		ttl,_:=rdb.TTL(rctx,cooldownKey).Result()
  return fmt.Errorf("please wait %d seconds before requesting again", int(ttl.Seconds()))



	}

	// increase the attempt counter 
	attempts, err := rdb.Incr(rctx, attemptKey).Result()
	if err != nil {
		return err
	}

	// expire after 1 hours
	rdb.Expire(rctx, attemptKey, time.Hour)


	var cooldown time.Duration

	
		switch attempts {
	case 1:
		cooldown = 60 * time.Second
	case 2:
		cooldown = 120 * time.Second
	case 3:
		cooldown = 240 * time.Second
	default:
		// Block for 1 hour

		rdb.Set(rctx, blockKey, "1", time.Hour)
		return fmt.Errorf("too many attempts. try again after 1 hour")
	}

	err = rdb.Set(rctx, cooldownKey, "1", cooldown).Err()
	if err != nil {
		return err
	}

	return nil





}