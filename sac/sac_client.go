package sac

import (
	"sac-cli/sac/dto"
)

type SecureAccessCloudClient interface {
	FindApplicationByName(name string) (*dto.ApplicationDTO, error)
	FindApplicationByID(id string) (*dto.ApplicationDTO, error)

	FindPolicyByName(name string) (dto.PolicyDTO, error)
	FindPoliciesByNames(name []string) ([]dto.PolicyDTO, error)
}
