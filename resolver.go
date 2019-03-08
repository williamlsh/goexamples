package goexamples

import (
	"context"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Video() VideoResolver {
	return &videoResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateVideo(ctx context.Context, input NewVideo) (*Video, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Videos(ctx context.Context, limit *int, offset *int) ([]Video, error) {
	panic("not implemented")
}

type videoResolver struct{ *Resolver }

func (r *videoResolver) CreatedAt(ctx context.Context, obj *Video) (string, error) {
	panic("not implemented")
}
