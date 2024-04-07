//go:build tools
// +build tools

package tools

import (
	_ "github.com/golang/mock/mockgen"
)

//go:generate mockgen -destination=internal/models/mocks/mock_contactrepo.go -package=mocks github.com/ethanjmarchand/learnhtmx/internal/models ContactRepo
//go:generate mockgen -destination=internal/models/mocks/mock_template.go -package=mocks github.com/ethanjmarchand/learnhtmx/internal/controllers Template
