package glance

import (
	"os"
	"testing"
)

func TestIsConfigStateValid(t *testing.T) {
	validPage := page{
		Title: "Home",
		Columns: []struct {
			Size    string  `yaml:"size"`
			Widgets widgets `yaml:"widgets"`
		}{
			{Size: "full"},
		},
	}

	t.Run("valid minimal config", func(t *testing.T) {
		cfg := &config{}
		cfg.Pages = []page{validPage}
		if err := isConfigStateValid(cfg); err != nil {
			t.Errorf("expected valid config, got error: %v", err)
		}
	})

	t.Run("no pages", func(t *testing.T) {
		cfg := &config{}
		if err := isConfigStateValid(cfg); err == nil {
			t.Error("expected error for no pages")
		}
	})

	t.Run("auth users without secret key", func(t *testing.T) {
		cfg := &config{}
		cfg.Pages = []page{validPage}
		cfg.Auth.Users = map[string]*user{
			"admin": {Password: "secret123"},
		}
		if err := isConfigStateValid(cfg); err == nil {
			t.Error("expected error when users are set without secret-key")
		}
	})

	t.Run("auth user with short username", func(t *testing.T) {
		cfg := &config{}
		cfg.Pages = []page{validPage}
		cfg.Auth.SecretKey = "some-key"
		cfg.Auth.Users = map[string]*user{
			"ab": {Password: "secret123"},
		}
		if err := isConfigStateValid(cfg); err == nil {
			t.Error("expected error for username shorter than 3 chars")
		}
	})

	t.Run("auth user with no password", func(t *testing.T) {
		cfg := &config{}
		cfg.Pages = []page{validPage}
		cfg.Auth.SecretKey = "some-key"
		cfg.Auth.Users = map[string]*user{
			"admin": {},
		}
		if err := isConfigStateValid(cfg); err == nil {
			t.Error("expected error for user with no password or hash")
		}
	})

	t.Run("auth user with short password", func(t *testing.T) {
		cfg := &config{}
		cfg.Pages = []page{validPage}
		cfg.Auth.SecretKey = "some-key"
		cfg.Auth.Users = map[string]*user{
			"admin": {Password: "abc"},
		}
		if err := isConfigStateValid(cfg); err == nil {
			t.Error("expected error for password shorter than 6 chars")
		}
	})

	t.Run("auth user with password hash instead of password", func(t *testing.T) {
		cfg := &config{}
		cfg.Pages = []page{validPage}
		cfg.Auth.SecretKey = "some-key"
		cfg.Auth.Users = map[string]*user{
			"admin": {PasswordHashString: "somehash"},
		}
		if err := isConfigStateValid(cfg); err != nil {
			t.Errorf("expected valid config with password hash, got error: %v", err)
		}
	})

	t.Run("invalid page width", func(t *testing.T) {
		cfg := &config{}
		p := validPage
		p.Width = "huge"
		cfg.Pages = []page{p}
		if err := isConfigStateValid(cfg); err == nil {
			t.Error("expected error for invalid page width")
		}
	})

	t.Run("slim page with too many columns", func(t *testing.T) {
		cfg := &config{}
		p := page{
			Title: "Home",
			Width: "slim",
			Columns: []struct {
				Size    string  `yaml:"size"`
				Widgets widgets `yaml:"widgets"`
			}{
				{Size: "full"},
				{Size: "small"},
				{Size: "small"},
			},
		}
		cfg.Pages = []page{p}
		if err := isConfigStateValid(cfg); err == nil {
			t.Error("expected error for slim page with more than 2 columns")
		}
	})

	t.Run("page with no full column", func(t *testing.T) {
		cfg := &config{}
		p := page{
			Title: "Home",
			Columns: []struct {
				Size    string  `yaml:"size"`
				Widgets widgets `yaml:"widgets"`
			}{
				{Size: "small"},
				{Size: "small"},
			},
		}
		cfg.Pages = []page{p}
		if err := isConfigStateValid(cfg); err == nil {
			t.Error("expected error for page with no full column")
		}
	})

	t.Run("page with no columns", func(t *testing.T) {
		cfg := &config{}
		p := page{Title: "Home"}
		cfg.Pages = []page{p}
		if err := isConfigStateValid(cfg); err == nil {
			t.Error("expected error for page with no columns")
		}
	})

	t.Run("invalid column size", func(t *testing.T) {
		cfg := &config{}
		p := page{
			Title: "Home",
			Columns: []struct {
				Size    string  `yaml:"size"`
				Widgets widgets `yaml:"widgets"`
			}{
				{Size: "medium"},
			},
		}
		cfg.Pages = []page{p}
		if err := isConfigStateValid(cfg); err == nil {
			t.Error("expected error for invalid column size")
		}
	})
}

func TestParseConfigVariables(t *testing.T) {
	t.Run("env variable substitution", func(t *testing.T) {
		os.Setenv("TEST_GLANCE_VAR", "my-value")
		defer os.Unsetenv("TEST_GLANCE_VAR")

		input := []byte("key: ${TEST_GLANCE_VAR}")
		result, err := parseConfigVariables(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(result) != "key: my-value" {
			t.Errorf("expected 'key: my-value', got %q", string(result))
		}
	})

	t.Run("missing env variable", func(t *testing.T) {
		os.Unsetenv("NONEXISTENT_GLANCE_VAR")
		input := []byte("key: ${NONEXISTENT_GLANCE_VAR}")
		_, err := parseConfigVariables(input)
		if err == nil {
			t.Error("expected error for missing env variable")
		}
	})

	t.Run("escaped variable", func(t *testing.T) {
		input := []byte(`key: \${NOT_REPLACED}`)
		result, err := parseConfigVariables(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(result) != "key: ${NOT_REPLACED}" {
			t.Errorf("expected escaped variable preserved, got %q", string(result))
		}
	})

	t.Run("no variables", func(t *testing.T) {
		input := []byte("key: plain-value")
		result, err := parseConfigVariables(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(result) != "key: plain-value" {
			t.Errorf("expected unchanged output, got %q", string(result))
		}
	})

	t.Run("lowercase variable name ignored", func(t *testing.T) {
		// Env variable names must match ^[A-Z0-9_]+$
		// Lowercase names should be passed through as-is
		input := []byte("key: ${lowercase_var}")
		result, err := parseConfigVariables(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(result) != "key: ${lowercase_var}" {
			t.Errorf("expected lowercase var to be passed through, got %q", string(result))
		}
	})

	t.Run("unknown variable type ignored", func(t *testing.T) {
		input := []byte("key: ${unknownType:something}")
		result, err := parseConfigVariables(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(result) != "key: ${unknownType:something}" {
			t.Errorf("expected unknown type to be passed through, got %q", string(result))
		}
	})
}

func TestOrderedYAMLMap(t *testing.T) {
	t.Run("new and iterate", func(t *testing.T) {
		keys := []string{"b", "a", "c"}
		values := []int{2, 1, 3}
		om, err := newOrderedYAMLMap(keys, values)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var gotKeys []string
		for k, _ := range om.Items() {
			gotKeys = append(gotKeys, k)
		}
		// Should preserve insertion order
		if len(gotKeys) != 3 || gotKeys[0] != "b" || gotKeys[1] != "a" || gotKeys[2] != "c" {
			t.Errorf("expected [b a c], got %v", gotKeys)
		}
	})

	t.Run("get", func(t *testing.T) {
		om, _ := newOrderedYAMLMap([]string{"x"}, []int{42})
		v, ok := om.Get("x")
		if !ok || v != 42 {
			t.Errorf("expected (42, true), got (%v, %v)", v, ok)
		}
		_, ok = om.Get("missing")
		if ok {
			t.Error("expected false for missing key")
		}
	})

	t.Run("merge preserves order", func(t *testing.T) {
		a, _ := newOrderedYAMLMap([]string{"x", "y"}, []int{1, 2})
		b, _ := newOrderedYAMLMap([]string{"y", "z"}, []int{20, 3})
		merged := a.Merge(b)

		var gotKeys []string
		for k, _ := range merged.Items() {
			gotKeys = append(gotKeys, k)
		}
		// x from a, y from a (overwritten by b), z from b
		if len(gotKeys) != 3 || gotKeys[0] != "x" || gotKeys[1] != "y" || gotKeys[2] != "z" {
			t.Errorf("expected [x y z], got %v", gotKeys)
		}
		// y should have b's value
		v, _ := merged.Get("y")
		if v != 20 {
			t.Errorf("expected y=20 after merge, got %v", v)
		}
	})

	t.Run("mismatched lengths", func(t *testing.T) {
		_, err := newOrderedYAMLMap([]string{"a", "b"}, []int{1})
		if err == nil {
			t.Error("expected error for mismatched key/value lengths")
		}
	})
}
