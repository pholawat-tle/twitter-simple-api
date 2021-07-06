package tweet

import "errors"

type Service struct{}

type TweetService interface {
	DeleteTweetByID(ID uint) error
	GetAllTweet() (Tweets, error)
	GetTweetByID(ID uint) (Tweet, error)
	CreateTweet(tweet Tweet) (Tweet, error)
}

type Tweet struct {
	Message string
	Author  string
	Likes   int
}

type Tweets []Tweet

var TweetDB = Tweets{
	Tweet{
		Message: "1st Tweet",
		Author:  "Pholawat Tangsatit",
		Likes:   0,
	},
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) DeleteTweetByID(ID uint) error {
	if ID > uint(len(TweetDB)-1) || len(TweetDB) == 0 {
		return errors.New("record doesn't exist")
	}
	TweetDB[len(TweetDB)-1], TweetDB[ID] = TweetDB[ID], TweetDB[len(TweetDB)-1]
	return nil
}

func (s *Service) GetAllTweet() (Tweets, error) {
	return TweetDB, nil
}

func (s *Service) GetTweetByID(ID uint) (Tweet, error) {
	if ID > uint(len(TweetDB)-1) || len(TweetDB) == 0 {
		return Tweet{}, errors.New("record doesn't exist")
	}
	return TweetDB[ID], nil
}

func (s *Service) CreateTweet(tweet Tweet) (Tweet, error) {
	TweetDB = append(TweetDB, tweet)
	return tweet, nil
}
