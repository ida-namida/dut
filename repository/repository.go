package repository

import (
	"dut/model"
)

type Repository interface {
	GetDecision(name string) (model.Decision, error)
}