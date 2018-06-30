package polluter

import (
	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/romanyx/jwalk"
)

type redisEngine struct {
	cli *redis.Client
}

func (e redisEngine) exec(cmds []command) error {
	for _, cmd := range cmds {
		if err := e.cli.Set(cmd.q, cmd.args[0], 0).Err(); err != nil {
			return errors.Wrap(err, "failed to set")
		}
	}
	return nil
}

func (e redisEngine) build(obj jwalk.ObjectWalker) (commands, error) {
	cmds := make(commands, 0)
	var err error

	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			}
		}
	}()

	obj.Walk(func(key string, value interface{}) {
		data, err := json.Marshal(value)
		if err != nil {
			panic(errors.Wrap(err, "marshal value"))
		}

		cmds = append(cmds, command{key, []interface{}{data}})
	})

	return cmds, err
}