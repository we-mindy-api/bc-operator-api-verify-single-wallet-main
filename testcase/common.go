package testcase

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"encoding/json"

	"github.com/PGITAb/bc-operator-api-verify/api"
)

func IsStatusCode(code int) api.TestCase {
	return api.TestCase{
		Name: "Expected status code " + fmt.Sprint(code),
		Func: func(rsp api.Response) error {
			if rsp.Code != code {
				return fmt.Errorf("status code incorrect, expect: %d, got: %d", code, rsp.Code)
			}
			return nil
		},
	}
}

func IsBodyMatchText(text string) api.TestCase {
	return api.TestCase{
		Name: "Expected body match text [" + text + "]",
		Func: func(rsp api.Response) error {
			if string(rsp.Body) != text {
				return fmt.Errorf("response do not match, expect: %s, got: %s", text, string(rsp.Body))
			}
			return nil
		},
	}
}

func ContainFields(keys []string) api.TestCase {
	return api.TestCase{
		Name: "Response json contain fields",
		Func: func(rsp api.Response) error {
			dict := make(map[string]interface{})
			err := json.Unmarshal([]byte(rsp.Body), &dict)
			if err != nil {
				return err
			}

			for _, k := range keys {
				key := k
				pt := dict
			Loop:
				for {
					if i := strings.Index(key, "/"); i != -1 {
						if _, ok := pt[key[0:i]]; !ok {
							return fmt.Errorf("response missing field: %s", k)
						}
						pt = pt[key[0:i]].(map[string]interface{})
						key = key[i+1:]
					} else {
						if _, ok := pt[key]; !ok {
							return fmt.Errorf("response missing field: %s", k)
						}
						break Loop
					}
				}
			}

			return nil
		},
	}
}

func DecodeToStruct(pointer interface{}) api.TestCase {
	return api.TestCase{
		Name: "Parse response json to struct",
		Func: func(rsp api.Response) error {
			err := json.Unmarshal(rsp.Body, pointer)
			if err != nil {
				return errors.New("parse json to struct failed")
			}
			return nil
		},
	}
}

func AssertIntField(key string, expected int) api.TestCase {
	return api.TestCase{
		Name: "Field value assertion",
		Func: func(rsp api.Response) error {
			// Parse json
			dict := make(map[string]interface{})
			err := json.Unmarshal([]byte(rsp.Body), &dict)
			if err != nil {
				return err
			}

			// Convert value
			actual := math.MinInt32
			if r, ok := dict[key]; ok {
				v, ok2 := r.(float64)
				if !ok2 {
					return fmt.Errorf("field \"%s\" is not integer", key)
				}
				actual = int(math.Round(v))
			} else {
				return fmt.Errorf("response missing field: %s", key)
			}

			// Assert value
			if actual != expected {
				return fmt.Errorf("\"%s\" value incorrect, expect: %d, got: %d", key, expected, actual)
			}

			return nil
		},
	}
}

func AssertStringField(key string, expected string) api.TestCase {
	return api.TestCase{
		Name: "Field value assertion",
		Func: func(rsp api.Response) error {
			// Parse json
			dict := make(map[string]interface{})
			err := json.Unmarshal([]byte(rsp.Body), &dict)
			if err != nil {
				return err
			}

			// Convert value
			actual := ""
			if r, ok := dict[key]; ok {
				actual, ok = r.(string)
				if !ok {
					return fmt.Errorf("field \"%s\" is not string", key)
				}
			} else {
				return fmt.Errorf("response missing field: %s", key)
			}

			// Assert value
			if actual != expected {
				return fmt.Errorf("\"%s\" value incorrect, expect: %s, got: %s", key, expected, actual)
			}

			return nil
		},
	}
}
