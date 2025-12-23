package subreddit

import (
	"time"

	"github.com/Andriy-Sydorenko/agora_backend/internal/user"
	"github.com/google/uuid"
)

type SubredditResponse struct {
	ID          uuid.UUID               `json:"id"`
	Name        string                  `json:"name"`
	DisplayName string                  `json:"display_name"`
	Description *string                 `json:"description,omitempty"`
	IconURL     *string                 `json:"icon_url,omitempty"`
	Creator     user.PublicUserResponse `json:"creator"`
	MemberCount int                     `json:"member_count"`
	PostCount   int                     `json:"post_count"`
	IsPublic    bool                    `json:"is_public"`
	IsNSFW      bool                    `json:"is_nsfw"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

type SubredditListResponse struct {
	Subreddits []SubredditResponse `json:"subreddits"`
}

type CreateSubredditRequest struct {
	Name        string  `json:"name"`
	DisplayName string  `json:"display_name"`
	Description *string `json:"description,omitempty"`
	IconURL     *string `json:"icon_url,omitempty"`
	IsPublic    *bool   `json:"is_public,omitempty"`
	IsNSFW      *bool   `json:"is_nsfw,omitempty"`
}

type UpdateSubredditRequest struct {
	DisplayName *string `json:"display_name"`
	Description *string `json:"description,omitempty"`
	IconURL     *string `json:"icon_url,omitempty"`
	IsPublic    *bool   `json:"is_public"`
	IsNSFW      *bool   `json:"is_nsfw"`
}

func ToSubredditResponse(s *Subreddit) SubredditResponse {
	return SubredditResponse{
		ID:          s.ID,
		Name:        s.Name,
		DisplayName: s.DisplayName,
		Description: s.Description,
		IconURL:     s.IconURL,
		Creator: user.PublicUserResponse{
			Username: s.Creator.Username,
			Email:    s.Creator.Email,
		},
		MemberCount: s.MemberCount,
		PostCount:   s.PostCount,
		IsPublic:    s.IsPublic,
		IsNSFW:      s.IsNSFW,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func ToSubredditListResponse(subreddits []Subreddit) SubredditListResponse {
	responses := make([]SubredditResponse, len(subreddits))
	for i := range subreddits {
		responses[i] = ToSubredditResponse(&subreddits[i])
	}
	return SubredditListResponse{
		Subreddits: responses,
	}
}
