package quill

import (
	"encoding/json"
	"fmt"
	"strings"
)

func Inspect(ops []byte, filters []InspectFilter) ([]InspectFilterResult, error) {
	return inspect(ops, filters)
}

func inspect(ops []byte, filters []InspectFilter) ([]InspectFilterResult, error) {

	var results []InspectFilterResult

	raw := make([]rawOp, 0, 12)
	if err := json.Unmarshal(ops, &raw); err != nil {
		return nil, fmt.Errorf("json.Unmarshal fail : %w", err)
	}

	for i := range raw {
		op := Op{Attrs: make(map[string]string, 3)}
		if err := raw[i].makeOp(&op); err != nil {
			return nil, fmt.Errorf("raw[i].makeOp(&op) fail : %w", err)
		}
		for _, filter := range filters {
			if filter.filtering(op) {
				results = append(results, InspectFilterResult{
					ByFilter: filter,
					Data:     op.Data,
				})
			}
		}
	}

	return results, nil
}

type InspectFilter struct {
	Type string
}

func (in InspectFilter) filtering(op Op) bool {

	if strings.TrimSpace(in.Type) != "" {
		if in.Type != op.Type {
			return false
		}
	}

	return true
}

type InspectFilterResult struct {
	ByFilter InspectFilter //ข้อมูลนี้ได้มาจาก filter ไหน
	Data     string
}
