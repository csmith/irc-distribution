package main

import "time"

const (
	actionPull   = "pull"
	actionPush   = "push"
	actionMount  = "mount"
	actionDelete = "delete"
)

type envelope struct {
	Events []event `json:"events"`
}

type event struct {
	Id        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Target    target    `json:"target"`
	Actor     actor     `json:"actor"`
}

type target struct {
	Digest     string `json:"digest"`
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	Size       int64  `json:"size"`
}

type actor struct {
	Name string `json:"name"`
}
