package iching

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/example/iching-fiber-app/internal/domain"
)

func RandomLine() (domain.LineValue, error) {
	sum := 0
	for i := 0; i < 3; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(2))
		if err != nil {
			return 0, fmt.Errorf("crypto rand coin toss: %w", err)
		}
		if n.Int64() == 0 {
			sum += 2
		} else {
			sum += 3
		}
	}
	return domain.LineValue(sum), nil
}

func RandomHexagramLines() ([]domain.Line, error) {
	lines := make([]domain.Line, 6)
	for i := 0; i < 6; i++ {
		v, err := RandomLine()
		if err != nil {
			return nil, err
		}
		lines[i] = domain.Line{Position: i + 1, Value: v}
	}
	return lines, nil
}
