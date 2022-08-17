package quill

import (
	"encoding/json"
	"fmt"
)

func Replace(ops []byte, cmds []ReplaceCmd) ([]byte, error) {
	return replace(ops, cmds)
}

func replace(ops []byte, cmds []ReplaceCmd) ([]byte, error) {

	raws := make([]rawOp, 0, 12)
	if err := json.Unmarshal(ops, &raws); err != nil {
		return nil, fmt.Errorf("json.Unmarshal fail : %w", err)
	}

	var rawOps []rawOp
	for i := range raws {
		op := Op{Attrs: make(map[string]string, 3)}
		if err := raws[i].makeOp(&op); err != nil {
			return nil, fmt.Errorf("raw[i].makeOp(&op) fail : %w", err)
		}
		for _, cmd := range cmds {
			if cmd.Replacer == nil {
				continue
			}
			if cmd.filtering(op) {
				err := cmd.Replacer(cmd, &op)
				if err != nil {
					return nil, fmt.Errorf("cmd.Replacer(&op) fail : %w", err)
				}
			}
		}

		err := raws[i].bind(op)
		if err != nil {
			return nil, fmt.Errorf("raws[i].bind(op) fail : %w", err)
		}
		rawOps = append(rawOps, raws[i])
	}

	buff, err := json.Marshal(rawOps)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal(buffOps) fail : %w", err)
	}

	//fmt.Printf("B= %s\n", string(buff))
	return buff, nil
}

type ReplaceCmd struct {
	InspectFilter
	Replacer (func(cmd ReplaceCmd, op *Op) error)
}
