package post

import "time"

type Post struct {
	Score            int       `json:"score" bson:"score"`
	Views            int       `json:"views" bson:"views"`
	Type             string    `json:"type" bson:"type"`
	Title            string    `json:"title" bson:"title"`
	Author           Author    `json:"author" bson:"author"`
	Category         string    `json:"category" bson:"category"`
	Text             string    `json:"text,omitempty" bson:"text"`
	Url              string    `json:"url,omitempty" bson:"url"`
	Votes            []Vote    `json:"votes" bson:"votes"`
	Comments         []Comment `json:"comments" bson:"comments"`
	Crated           time.Time `json:"created" bson:"created"`
	UpvotePercentage int       `json:"upvotePercentage" bson:"upvotePercentage"`
	ID               string    `json:"id" bson:"_id"`
}

type Comment struct {
	Crated time.Time `json:"created" bson:"created"`
	Author Author    `json:"author" bson:"author"`
	Body   string    `json:"body" bson:"body"`
	ID     string    `json:"id" bson:"id"`
}

type Vote struct {
	User string `json:"user" bson:"user"`
	Vote int    `json:"vote" bson:"vote"`
}

type Author struct {
	UserName string `json:"username" bson:"username"`
	ID       string `json:"id" bson:"id"`
}

func (post *Post) CalculateUpvotePercentage() {
	if len(post.Votes) == 0 {
		post.UpvotePercentage = 0
		return
	}
	plusVotes := 0
	for _, nextVote := range post.Votes {
		if nextVote.Vote == 1 {
			plusVotes++
		}
	}
	post.UpvotePercentage = (plusVotes * 100) / len(post.Votes)
}

func (post *Post) CalculateScore() {
	score := 0
	for _, nextVote := range post.Votes {
		score = score + nextVote.Vote
	}
	post.Score = score
}

func (post *Post) AddComment(comment Comment) {
	post.Comments = append(post.Comments, comment)
}

func (post *Post) DeleteComment(id string) {
	for i, nextComment := range post.Comments {
		if nextComment.ID == id {
			post.Comments = append(post.Comments[:i], post.Comments[i+1:]...)
			return
		}
	}
}

func (post *Post) Upvote(userId string) {
	post.Unvote(userId)
	newVote := Vote{
		User: userId,
		Vote: 1,
	}
	post.Votes = append(post.Votes, newVote)
	post.CalculateUpvotePercentage()
	post.CalculateScore()
}

func (post *Post) Downvote(userId string) {
	post.Unvote(userId)
	newVote := Vote{
		User: userId,
		Vote: -1,
	}
	post.Votes = append(post.Votes, newVote)
	post.CalculateUpvotePercentage()
	post.CalculateScore()
}

func (post *Post) Unvote(userId string) {
	for i, nextVote := range post.Votes {
		if nextVote.User == userId {
			post.Votes = append(post.Votes[:i], post.Votes[i+1:]...)
			post.CalculateUpvotePercentage()
			post.CalculateScore()
			return
		}
	}
}
