package main

import (
	"fmt"
	"time"
)

func fixtureBlogCreate(i int) *Blog {
	return &Blog{
		ID:        1 * i,
		Title:     fmt.Sprintf("title %d", i),
		CreatedAt: time.Now(),
		Posts: []*Post{
			{
				ID:    1 * i,
				Title: "Foo",
				Body:  "Bar",
				Comments: []*Comment{
					{
						ID:   1 * i,
						Body: "Foo",
					},
					{
						ID:   2 * i,
						Body: "Bar",
					},
				},
			},
			{
				ID:    2 * i,
				Title: "Foo",
				Body:  "Bar",
				Comments: []*Comment{
					{
						ID:   1 * i,
						Body: "Foo",
					},
					{
						ID:   3 * i,
						Body: "Bar",
					},
				},
			},
		},
		CurrentPost: &Post{
			ID:    1 * i,
			Title: "Foo",
			Body:  "Bar",
			Comments: []*Comment{
				{
					ID:   1 * i,
					Body: "foo",
				},
				{
					ID:   2 * i,
					Body: "bar",
				},
			},
		},
	}
}

func fixtureBlogList() (blogs []interface{}) {
	for i := 1; i < 2; i++ {
		blogs = append(blogs, fixtureBlogCreate(i))
	}
	return
}
