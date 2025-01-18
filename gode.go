package gode

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

const runnerNode = `(function(program, execJS) { execJS(program) })(function() { #{source} 
}, function(program) {
  var output;
  var print = function(string) {
    process.stdout.write('' + string + '\n');
  };
  try {
    result = program();
    print('')
    if (typeof result == 'undefined' && result !== null) {
      print('["ok"]');
    } else {
      try {
        print(JSON.stringify(['ok', result]));
      } catch (err) {
        print('["err"]');
      }
    }
  } catch (err) {
    print(JSON.stringify(['err', '' + err]));
  }
});`

type Gode struct {
	ctx    context.Context
	source string
	dir    string
}

// New returns a new Gode instance.
func New(source ...string) (*Gode, error) {
	return NewWithContext(context.Background(), source...)
}

// NewWithContext returns a new Gode instance with context.
func NewWithContext(ctx context.Context, source ...string) (*Gode, error) {
	var src string
	if len(source) > 1 {
		return nil, errors.New("invalid parameter")
	}
	if len(source) > 0 {
		src = source[0]
	}

	gode := &Gode{ctx: ctx, source: src}

	if err := gode.isAvailable(); err != nil {
		return nil, err
	}
	return gode, nil
}

// WorkPath sets the working directory for the Gode instance.
func (g *Gode) WorkPath(dir string) {
	g.dir = dir
}

// Eval evaluates the given source code.
// It handles empty source by setting data to "”", otherwise it formats the source code.
func (g *Gode) Eval(source string) (interface{}, error) {
	var data string
	if strings.TrimSpace(source) == "" {
		data = "''"
	} else {
		data = fmt.Sprintf("'('+%s+')'", fmt.Sprintf("%q", source))
	}

	code := fmt.Sprintf("return eval(%s)", data)

	return g.exec(code)
}

// Call calls a function with the provided arguments.
// It marshals the arguments into JSON format and constructs the code for evaluation.
func (g *Gode) Call(function string, args ...interface{}) (interface{}, error) {
	argsJson, _ := json.Marshal(args)
	code := fmt.Sprintf("%s.apply(this, %s)", function, argsJson)
	return g.Eval(code)
}

// exec executes the provided source code.
// It combines the Gode's source with the new source if necessary.
// It uses exec.CommandContext to run the code in a node environment.
func (g *Gode) exec(source string) (interface{}, error) {
	if g.source != "" {
		source = g.source + "\n" + source
	}

	input := strings.ReplaceAll(runnerNode, "#{source}", source)

	cmd := exec.CommandContext(g.ctx, "node")
	cmd.Stdin = strings.NewReader(input)
	if g.dir != "" {
		cmd.Dir = g.dir
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("output:%s err:%+v", output, err)
	}
	return g.result(output)
}

// result processes the output from the execution.
// It removes carriage returns, splits the output by new lines, and unmarshals the last but one line.
// It checks the status and returns an error if it's not "ok".
func (g *Gode) result(output []byte) (interface{}, error) {
	output = bytes.ReplaceAll(output, []byte{'\r'}, []byte{})
	lines := bytes.Split(output, []byte{'\n'})

	data := lines[len(lines)-2]

	var res []interface{}
	if err := json.Unmarshal(data, &res); err != nil {
		return "", fmt.Errorf("json unmarshal data:%s err:%+v", data, err)
	}

	if len(res) == 0 {
		res = append(res, "err")
	}
	if len(res) == 1 {
		res = append(res, nil)
	}

	status, value := res[0].(string), res[1]
	if status != "ok" {
		return "", fmt.Errorf("err:%+v", value)
	}
	return value, nil
}

// isAvailable checks if the node command is available.
func (g *Gode) isAvailable() error {
	cmd := exec.CommandContext(g.ctx, "node", "-v")
	return cmd.Run()
}
