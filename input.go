package gocommend

import "log"

type Input struct {
	Collection string
	UserId     string
	ItemId     string
	Rate       int
}

func Import(i *Input) error {
	log.Println(i)
	cSet := collectionSet{}
	cSet.init(i.Collection)

	if i.Rate > 0 {
		log.Println("like")
		if err := like(&cSet, i.UserId, i.ItemId); err != nil {
			return err
		}
	} else {
		log.Println("dislike")
		if err := dislike(&cSet, i.UserId, i.ItemId); err != nil {
			return err
		}
	}
	return nil
}

func Update(i *Input) error {

	log.Println(i)
	if i.Collection == "" {
		return gocommendError{emptyCollection}
	}

	algo := algorithms{}
	algo.cSet.init(i.Collection)

	// update specific user's sets
	if i.UserId != "" {
		if err := algo.updateSimilarityFor(i.UserId); err != nil {
			return err
		}
		if err := algo.updateRecommendationFor(i.UserId); err != nil {
			return err
		}
	}
	if i.ItemId != "" {
		if err := algo.updateWilsonScore(i.ItemId); err != nil {
			return err
		}
	}

	return nil
}

func like(cSet *collectionSet, userId string, itemId string) error {
	var (
		rs  interface{}
		err error
	)
	if rs, err = redisClient.Do("SISMEMBER", cSet.itemLiked(itemId), userId); err != nil {
		return err
	}
	if sis, _ := rs.(int); sis == 0 {
		redisClient.Do("ZINCRBY", cSet.mostLiked, 1, itemId)
	}
	if _, err = redisClient.Do("SADD", cSet.userLiked(userId), itemId); err != nil {
		return err
	}
	if _, err = redisClient.Do("SADD", cSet.itemLiked(itemId), userId); err != nil {
		return err
	}

	return nil
}

func dislike(cSet *collectionSet, userId string, itemId string) error {
	var (
		rs  interface{}
		err error
	)
	if rs, err = redisClient.Do("SISMEMBER", cSet.itemDisliked(itemId), userId); err != nil {
		return err
	}
	if sis, _ := rs.(int); sis == 0 {
		redisClient.Do("ZINCRBY", cSet.mostDisliked, 1, itemId)
	}
	if _, err = redisClient.Do("SADD", cSet.userDisliked(userId), itemId); err != nil {
		return err
	}
	if _, err = redisClient.Do("SADD", cSet.itemDisliked(itemId), userId); err != nil {
		return err
	}

	return nil
}
